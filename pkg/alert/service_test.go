package alert

import (
	"errors"
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/pkg/alert/alertmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"reflect"
	"testing"
)

type AlertmanagerClientMock struct {
	mock.Mock
}

func (am *AlertmanagerClientMock) SyncConfig(config alertmanager.AlertManagerConfig) error {
	args := am.Called(config)
	return args.Error(0)
}

func TestServiceUpsert(t *testing.T) {
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	var dsn string
	if postgresPassword == "" {
		dsn = "host=localhost user=postgres dbname=postgres port=5432 sslmode=disable"
	} else {
		dsn = fmt.Sprintf("host=localhost password=%s user=postgres dbname=postgres port=5432 sslmode=disable", postgresPassword)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	expectedConfigs := alertmanager.AlertManagerConfig{
		AlertHistoryHost: "http://example.com",
		EntityCredentials: alertmanager.EntityCredentials{
			Entity: "avengers",
			Token:  "random-token",
			Teams: map[string]alertmanager.TeamCredentials{
				"hydra": {
					Name:                "hydra",
					PagerdutyCredential: "abc",
					Slackcredentials: alertmanager.SlackConfig{
						Critical: alertmanager.SlackCredential{
							Channel: "critical_channel",
						},
						Warning: alertmanager.SlackCredential{
							Channel: "warning_channel2",
						},
					},
				},
				"wakanda": {
					Name:                "wakanda",
					PagerdutyCredential: "xyzw",
					Slackcredentials: alertmanager.SlackConfig{
						Critical: alertmanager.SlackCredential{
							Channel: "critical_channel",
						},
						Warning: alertmanager.SlackCredential{
							Channel: "warning_channel",
						},
					},
				},
			},
		},
	}

	t.Run("should insert new records", func(t *testing.T) {
		db.Exec("truncate slack_credentials, pagerduty_credentials")
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(SlackCredential{})
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(PagerdutyCredential{})
		if err != nil {
			t.Fatal(err)
		}
		clientMock := AlertmanagerClientMock{}
		clientMock.On("SyncConfig", mock.Anything).Return(nil)
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		mockedCodeExchangeService.On("GetToken", "avengers").
			Return("random-token", nil).Once()
		service := NewService(db, &clientMock, domain.SirenServiceConfig{
			Host: "http://example.com",
		}, mockedCodeExchangeService)

		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "critical_channel",
				},
				Warning: domain.SlackCredential{
					Channel: "warning_channel",
				},
			},
		}
		err = service.Upsert(credential)
		if err != nil {
			t.Fatal(err)
		}
		var warningSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "WARNING").First(&warningSlackCredential)
		assert.Equal(t, "warning_channel", warningSlackCredential.Channel)

		var criticalSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "CRITICAL").First(&criticalSlackCredential)
		assert.Equal(t, "critical_channel", criticalSlackCredential.Channel)
		var pagerdutyCredential PagerdutyCredential
		db.Model(&PagerdutyCredential{}).Where("team_name = ?", "hydra").First(&pagerdutyCredential)
		assert.Equal(t, "xyz", pagerdutyCredential.ServiceKey)
	})

	t.Run("should handle errors in getting slack token for given entity", func(t *testing.T) {
		clientMock := AlertmanagerClientMock{}
		clientMock.On("SyncConfig", mock.Anything).Return(nil).Once()
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		mockedCodeExchangeService.On("GetToken", "avengers").
			Return("", errors.New("random error")).Once()
		service := NewService(db, &clientMock, domain.SirenServiceConfig{
			Host: "http://example.com",
		}, mockedCodeExchangeService)

		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "critical_channel",
				},
				Warning: domain.SlackCredential{
					Channel: "warning_channel",
				},
			},
		}
		err := service.Upsert(credential)
		assert.EqualError(t, err, "failed to get token for entity avengers: random error")
	})

	t.Run("should upsert records", func(t *testing.T) {
		db.Exec("truncate slack_credentials, pagerduty_credentials")
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(SlackCredential{})
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(PagerdutyCredential{})
		if err != nil {
			t.Fatal(err)
		}
		clientMock := AlertmanagerClientMock{}
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		mockedCodeExchangeService.On("GetToken", "avengers").
			Return("random-token", nil).Once()
		clientMock.On("SyncConfig", mock.Anything).Return(nil)
		service := NewService(db, &clientMock, domain.SirenServiceConfig{
			Host: "http://example.com",
		}, mockedCodeExchangeService)
		result := db.Model(SlackCredential{}).Create(&SlackCredential{
			Channel:  "critical_channel",
			Level:    "CRITICAL",
			TeamName: "hydra",
			Entity:   "avengers",
		})
		result = db.Model(SlackCredential{}).Create(&SlackCredential{
			Channel:  "warning_channel",
			Level:    "WARNING",
			TeamName: "hydra",
			Entity:   "avengers",
		})
		db.Model(PagerdutyCredential{}).Create(&PagerdutyCredential{
			ServiceKey: "xyz",
			TeamName:   "hydra",
			Entity:     "avengers",
		})
		assert.Nil(t, result.Error)
		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "abc",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "critical_channel",
				},
				Warning: domain.SlackCredential{
					Channel: "warning_channel2",
				},
			},
		}
		err = service.Upsert(credential)
		if err != nil {
			t.Fatal(err)
		}
		var warningSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "WARNING").First(&warningSlackCredential)
		assert.Equal(t, "warning_channel2", warningSlackCredential.Channel)

		var criticalSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "CRITICAL").First(&criticalSlackCredential)
		assert.Equal(t, "critical_channel", criticalSlackCredential.Channel)
		var pagerdutyCredential PagerdutyCredential
		db.Model(&PagerdutyCredential{}).Where("team_name = ?", "hydra").First(&pagerdutyCredential)
		assert.Equal(t, "abc", pagerdutyCredential.ServiceKey)
	})

	t.Run("should update entity config for the alertmanager", func(t *testing.T) {
		db.Exec("truncate slack_credentials, pagerduty_credentials")
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(SlackCredential{})
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(PagerdutyCredential{})
		if err != nil {
			t.Fatal(err)
		}
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		mockedCodeExchangeService.On("GetToken", "avengers").
			Return("random-token", nil).Once()
		clientMock := AlertmanagerClientMock{}
		clientMock.On("SyncConfig", mock.MatchedBy(func(actualConfig alertmanager.AlertManagerConfig) bool {

			return reflect.DeepEqual(actualConfig, expectedConfigs)
		})).Return(nil)
		service := NewService(db, &clientMock, domain.SirenServiceConfig{
			Host: "http://example.com",
		}, mockedCodeExchangeService)
		result := db.Model(SlackCredential{}).Create(&SlackCredential{
			Channel:  "critical_channel",
			Level:    "CRITICAL",
			TeamName: "wakanda",
			Entity:   "avengers",
		})
		result = db.Model(SlackCredential{}).Create(&SlackCredential{
			Channel:  "warning_channel",
			Level:    "WARNING",
			TeamName: "wakanda",
			Entity:   "avengers",
		})
		db.Model(PagerdutyCredential{}).Create(&PagerdutyCredential{
			ServiceKey: "xyzw",
			TeamName:   "wakanda",
			Entity:     "avengers",
		})
		assert.Nil(t, result.Error)
		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "abc",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "critical_channel",
				},
				Warning: domain.SlackCredential{
					Channel: "warning_channel2",
				},
			},
		}
		err = service.Upsert(credential)
		if err != nil {
			t.Fatal(err)
		}
		var warningSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "WARNING").First(&warningSlackCredential)
		assert.Equal(t, "warning_channel2", warningSlackCredential.Channel)

		var criticalSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "CRITICAL").First(&criticalSlackCredential)
		assert.Equal(t, "critical_channel", criticalSlackCredential.Channel)
		var pagerdutyCredential PagerdutyCredential
		db.Model(&PagerdutyCredential{}).Where("team_name = ?", "hydra").First(&pagerdutyCredential)
		assert.Equal(t, "abc", pagerdutyCredential.ServiceKey)
		clientMock.AssertExpectations(t)
		clientMock.AssertCalled(t, "SyncConfig", expectedConfigs)

	})
}


func TestServiceGet(t *testing.T) {
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	var dsn string
	if postgresPassword == "" {
		dsn = "host=localhost user=postgres dbname=postgres port=5432 sslmode=disable"
	} else {
		dsn = fmt.Sprintf("host=localhost password=%s user=postgres dbname=postgres port=5432 sslmode=disable", postgresPassword)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should return alert credentials of the team", func(t *testing.T) {
		db.Exec("truncate slack_credentials, pagerduty_credentials")
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(SlackCredential{})
		if err != nil {
			t.Fatal(err)
		}
		err = db.AutoMigrate(PagerdutyCredential{})
		if err != nil {
			t.Fatal(err)
		}

		result := db.Model(SlackCredential{}).Create(&SlackCredential{
			Channel:  "critical_channel",
			Level:    "CRITICAL",
			TeamName: "hydra",
			Entity:   "avengers",
		})
		assert.Nil(t, result.Error)
		result = db.Model(SlackCredential{}).Create(&SlackCredential{
			Channel:  "warning_channel",
			Level:    "WARNING",
			TeamName: "hydra",
			Entity:   "avengers",
		})
		assert.Nil(t, result.Error)
		db.Model(PagerdutyCredential{}).Create(&PagerdutyCredential{
			ServiceKey: "xyz",
			TeamName:   "hydra",
			Entity:     "avengers",
		})
		assert.Nil(t, result.Error)
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		mockedCodeExchangeService.On("GetToken", "avengers").
			Return("random-token", nil).Once()
		service := NewService(db, nil, domain.SirenServiceConfig{
			Host: "http://example.com",
		}, mockedCodeExchangeService)
		credential, err := service.Get("hydra")
		assert.Nil(t, err)
		expectedCredential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel:  "critical_channel",
				},
				Warning: domain.SlackCredential{
					Channel:  "warning_channel",
				},
			},
		}
		assert.Equal(t, expectedCredential, credential)
	})
}
