package alert

import (
	"database/sql"
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/alert/alertmanager"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	repository          *Repository
	codeexchangeService domain.CodeExchangeService
	amClient            alertmanager.Client
	sirenHost           string
}

type SlackCredential struct {
	gorm.Model
	ChannelName string
	Level       string `gorm:"uniqueIndex:team_level_unique"`

	TeamName string `gorm:"uniqueIndex:team_level_unique"`
	Entity   string
}

type PagerdutyCredential struct {
	gorm.Model
	ServiceKey string
	TeamName   string `gorm:"uniqueIndex:team_unique"`
	Entity     string
}

func (s Service) Migrate() error {
	err := s.repository.db.AutoMigrate(SlackCredential{}, PagerdutyCredential{})
	if err != nil {
		return err
	}
	err = s.repository.db.Migrator().DropColumn(&SlackCredential{}, "username")
	if err != nil {
		return err
	}
	return s.repository.db.Migrator().DropColumn(&SlackCredential{}, "webhook")
}

func (s Service) Upsert(credential domain.AlertCredential) error {
	token, err := s.codeexchangeService.GetToken(credential.Entity)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get token for entity %s", credential.Entity))
	}
	err = s.repository.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(upsertCriticalSlackCredentialQuery,
			sql.Named("channel_name", credential.SlackConfig.Critical.Channel),
			sql.Named("entity", credential.Entity),
			sql.Named("team_name", credential.TeamName),
		).Error; err != nil {
			return err
		}
		if err := tx.Exec(upsertWarningSlackCredentialQuery,
			sql.Named("channel_name", credential.SlackConfig.Warning.Channel),
			sql.Named("entity", credential.Entity),
			sql.Named("team_name", credential.TeamName),
		).Error; err != nil {
			return err
		}

		if err := tx.Exec(upsertPagerdutyCredentialsQuery,
			sql.Named("entity", credential.Entity),
			sql.Named("service_key", credential.PagerdutyCredentials),
			sql.Named("team_name", credential.TeamName),
		).Error
			err != nil {
			return err
		}

		var credentials []SlackCredential
		tx.Model(SlackCredential{}).Where("entity= ?", credential.Entity).Find(&credentials)
		rows, err := tx.Raw(joinQuery, credential.Entity, credential.Entity, credential.Entity).Rows()
		defer rows.Close()
		teamCredentials := make(map[string]alertmanager.TeamCredentials)
		for rows.Next() {
			teamCredential := alertmanager.TeamCredentials{
				Slackcredentials: alertmanager.SlackConfig{
					Critical: alertmanager.SlackCredential{},
					Warning:  alertmanager.SlackCredential{},
				},
			}

			rows.Scan(&teamCredential.Name, &teamCredential.PagerdutyCredential,
				&teamCredential.Slackcredentials.Warning.Channel,
				&teamCredential.Slackcredentials.Critical.Channel)
			teamCredentials[teamCredential.Name] = teamCredential

		}
		err = s.amClient.SyncConfig(alertmanager.AlertManagerConfig{
			AlertHistoryHost: s.sirenHost,
			EntityCredentials: alertmanager.EntityCredentials{
				Entity: credential.Entity,
				Teams:  teamCredentials,
				Token:  token,
			},
		})

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) Get(teamName string) (domain.AlertCredential, error) {
	slackCredentialMap, pagerdutyCredential, err := s.repository.GetCredential(teamName)
	if err != nil {
		return domain.AlertCredential{}, err
	}
	warningSlackCredential := slackCredentialMap["WARNING"]
	criticalSlackCredential := slackCredentialMap["CRITICAL"]
	credential := domain.AlertCredential{
		Entity:               pagerdutyCredential.Entity,
		TeamName:             teamName,
		PagerdutyCredentials: pagerdutyCredential.ServiceKey,
		SlackConfig: domain.SlackConfig{
			Warning: domain.SlackCredential{
				Channel: warningSlackCredential.ChannelName,
			},
			Critical: domain.SlackCredential{
				Channel: criticalSlackCredential.ChannelName,
			},
		},
	}
	return credential, nil
}

func NewService(db *gorm.DB, amClient alertmanager.Client, c domain.SirenServiceConfig,
	codeExchangeSvc domain.CodeExchangeService) domain.AlertmanagerService {
	return &Service{
		repository:          NewRepository(db),
		amClient:            amClient,
		sirenHost:           c.Host,
		codeexchangeService: codeExchangeSvc,
	}
}
