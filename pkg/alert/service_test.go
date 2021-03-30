package alert

import (
	"fmt"
	"github.com/odpf/siren/domain"
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

func (am *AlertmanagerClientMock) SyncConfig(credentials alertmanager.EntityCredentials) error {
	args := am.Called(credentials)
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
	expectedEntityCredentials := alertmanager.EntityCredentials{
		Entity: "avengers",
		Teams: map[string]alertmanager.TeamCredentials{
			"hydra": {
				Name:                "hydra",
				PagerdutyCredential: "abc",
				Slackcredentials: alertmanager.SlackConfig{
					Critical: alertmanager.SlackCredential{
						Webhook:  "http://critical.com",
						Channel:  "critical_channel",
						Username: "critical_user",
					},
					Warning: alertmanager.SlackCredential{
						Webhook:  "http://warning.com",
						Channel:  "warning_channel2",
						Username: "warning_user",
					},
				},
			},
			"wakanda": {
				Name:                "wakanda",
				PagerdutyCredential: "xyzw",
				Slackcredentials: alertmanager.SlackConfig{
					Critical: alertmanager.SlackCredential{
						Webhook:  "http://criticalwakanda.com",
						Channel:  "critical_channel",
						Username: "critical_user",
					},
					Warning: alertmanager.SlackCredential{
						Webhook:  "http://warningwakanda.com",
						Channel:  "warning_channel",
						Username: "warning_user",
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
		service := NewService(db, &clientMock)

		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel:  "critical_channel",
					Webhook:  "http://critical.com",
					Username: "critical_user",
				},
				Warning: domain.SlackCredential{
					Channel:  "warning_channel",
					Webhook:  "http://warning.com",
					Username: "warning_user",
				},
			},
		}
		err = service.Upsert(credential)
		if err != nil {
			t.Fatal(err)
		}
		var warningSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "WARNING").First(&warningSlackCredential)
		assert.Equal(t, "http://warning.com", warningSlackCredential.Webhook)
		assert.Equal(t, "warning_user", warningSlackCredential.Username)
		assert.Equal(t, "warning_channel", warningSlackCredential.ChannelName)

		var criticalSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "CRITICAL").First(&criticalSlackCredential)
		assert.Equal(t, "http://critical.com", criticalSlackCredential.Webhook)
		assert.Equal(t, "critical_user", criticalSlackCredential.Username)
		assert.Equal(t, "critical_channel", criticalSlackCredential.ChannelName)
		var pagerdutyCredential PagerdutyCredential
		db.Model(&PagerdutyCredential{}).Where("team_name = ?", "hydra").First(&pagerdutyCredential)
		assert.Equal(t, "xyz", pagerdutyCredential.ServiceKey)

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
		clientMock.On("SyncConfig", mock.Anything).Return(nil)
		service := NewService(db, &clientMock)
		result := db.Model(SlackCredential{}).Create(&SlackCredential{
			ChannelName: "critical_channel",
			Username:    "critical_user",
			Webhook:     "http://critical.com",
			Level:       "CRITICAL",
			TeamName:    "hydra",
			Entity:      "avengers",
		})
		result = db.Model(SlackCredential{}).Create(&SlackCredential{
			ChannelName: "warning_channel",
			Username:    "warning_user",
			Webhook:     "http://warning.com",
			Level:       "WARNING",
			TeamName:    "hydra",
			Entity:      "avengers",
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
					Channel:  "critical_channel",
					Webhook:  "http://critical.com",
					Username: "critical_user",
				},
				Warning: domain.SlackCredential{
					Channel:  "warning_channel2",
					Webhook:  "http://warning.com",
					Username: "warning_user",
				},
			},
		}
		err = service.Upsert(credential)
		if err != nil {
			t.Fatal(err)
		}
		var warningSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "WARNING").First(&warningSlackCredential)
		assert.Equal(t, "http://warning.com", warningSlackCredential.Webhook)
		assert.Equal(t, "warning_user", warningSlackCredential.Username)
		assert.Equal(t, "warning_channel2", warningSlackCredential.ChannelName)

		var criticalSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "CRITICAL").First(&criticalSlackCredential)
		assert.Equal(t, "http://critical.com", criticalSlackCredential.Webhook)
		assert.Equal(t, "critical_user", criticalSlackCredential.Username)
		assert.Equal(t, "critical_channel", criticalSlackCredential.ChannelName)
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
		clientMock := AlertmanagerClientMock{}
		clientMock.On("SyncConfig", mock.MatchedBy(func(actualCredentials alertmanager.EntityCredentials) bool {

			return reflect.DeepEqual(actualCredentials, expectedEntityCredentials)
		})).Return(nil)
		service := NewService(db, &clientMock)
		result := db.Model(SlackCredential{}).Create(&SlackCredential{
			ChannelName: "critical_channel",
			Username:    "critical_user",
			Webhook:     "http://criticalwakanda.com",
			Level:       "CRITICAL",
			TeamName:    "wakanda",
			Entity:      "avengers",
		})
		result = db.Model(SlackCredential{}).Create(&SlackCredential{
			ChannelName: "warning_channel",
			Username:    "warning_user",
			Webhook:     "http://warningwakanda.com",
			Level:       "WARNING",
			TeamName:    "wakanda",
			Entity:      "avengers",
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
					Channel:  "critical_channel",
					Webhook:  "http://critical.com",
					Username: "critical_user",
				},
				Warning: domain.SlackCredential{
					Channel:  "warning_channel2",
					Webhook:  "http://warning.com",
					Username: "warning_user",
				},
			},
		}
		err = service.Upsert(credential)
		if err != nil {
			t.Fatal(err)
		}
		var warningSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "WARNING").First(&warningSlackCredential)
		assert.Equal(t, "http://warning.com", warningSlackCredential.Webhook)
		assert.Equal(t, "warning_user", warningSlackCredential.Username)
		assert.Equal(t, "warning_channel2", warningSlackCredential.ChannelName)

		var criticalSlackCredential SlackCredential
		db.Model(&SlackCredential{}).Where("team_name = ? AND level =?", "hydra", "CRITICAL").First(&criticalSlackCredential)
		assert.Equal(t, "http://critical.com", criticalSlackCredential.Webhook)
		assert.Equal(t, "critical_user", criticalSlackCredential.Username)
		assert.Equal(t, "critical_channel", criticalSlackCredential.ChannelName)
		var pagerdutyCredential PagerdutyCredential
		db.Model(&PagerdutyCredential{}).Where("team_name = ?", "hydra").First(&pagerdutyCredential)
		assert.Equal(t, "abc", pagerdutyCredential.ServiceKey)
		clientMock.AssertExpectations(t)
		clientMock.AssertCalled(t, "SyncConfig", expectedEntityCredentials)

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
			ChannelName: "critical_channel",
			Username:    "critical_user",
			Webhook:     "http://critical.com",
			Level:       "CRITICAL",
			TeamName:    "hydra",
			Entity:      "avengers",
		})
		assert.Nil(t, result.Error)
		result = db.Model(SlackCredential{}).Create(&SlackCredential{
			ChannelName: "warning_channel",
			Username:    "warning_user",
			Webhook:     "http://warning.com",
			Level:       "WARNING",
			TeamName:    "hydra",
			Entity:      "avengers",
		})
		assert.Nil(t, result.Error)
		db.Model(PagerdutyCredential{}).Create(&PagerdutyCredential{
			ServiceKey: "xyz",
			TeamName:   "hydra",
			Entity:     "avengers",
		})
		assert.Nil(t, result.Error)
		service := NewService(db, nil)
		credential, err := service.Get("hydra")
		assert.Nil(t, err)
		expectedCredential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel:  "critical_channel",
					Webhook:  "http://critical.com",
					Username: "critical_user",
				},
				Warning: domain.SlackCredential{
					Channel:  "warning_channel",
					Webhook:  "http://warning.com",
					Username: "warning_user",
				},
			},
		}
		assert.Equal(t, expectedCredential, credential)
	})
}
