package alert

import (
	"database/sql"
	"github.com/odpf/siren/alert/alertmanager"
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
)

type Service struct {
	repository *Repository
	amClient   alertmanager.Client
}

type SlackCredential struct {
	gorm.Model
	ChannelName string
	Username    string
	Webhook     string
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
	return s.repository.db.AutoMigrate(SlackCredential{}, PagerdutyCredential{})
}

func (s Service) Upsert(credential domain.AlertCredential) error {
	err := s.repository.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`INSERT INTO slack_credentials (created_at,updated_at,channel_name,username,webhook,level,team_name,entity)
  VALUES (now(), now(), @channel_name, @username, @webhook,'CRITICAL', @team_name,@entity) 
  ON CONFLICT (level, team_name)   
  DO UPDATE SET "updated_at"= now(),"deleted_at"="excluded"."deleted_at","channel_name"="excluded"."channel_name","username"="excluded"."username","webhook"="excluded"."webhook","level"="excluded"."level","team_name"="excluded"."team_name","entity"="excluded"."entity" RETURNING "id"`,
			sql.Named("channel_name", credential.SlackConfig.Critical.Channel),
			sql.Named("username", credential.SlackConfig.Critical.Username),
			sql.Named("webhook", credential.SlackConfig.Critical.Webhook),
			sql.Named("entity", credential.Entity),
			sql.Named("team_name", credential.TeamName),
		).Error; err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO slack_credentials (created_at,updated_at,channel_name,username,webhook,level,team_name,entity)
  VALUES (now(), now(), @channel_name, @username, @webhook,'WARNING', @team_name,@entity) 
  ON CONFLICT (level, team_name)   
  DO UPDATE SET "updated_at"= now(),"deleted_at"="excluded"."deleted_at","channel_name"="excluded"."channel_name","username"="excluded"."username","webhook"="excluded"."webhook","level"="excluded"."level","team_name"="excluded"."team_name","entity"="excluded"."entity" RETURNING "id"`,
			sql.Named("channel_name", credential.SlackConfig.Warning.Channel),
			sql.Named("username", credential.SlackConfig.Warning.Username),
			sql.Named("webhook", credential.SlackConfig.Warning.Webhook),
			sql.Named("entity", credential.Entity),
			sql.Named("team_name", credential.TeamName),
		).Error; err != nil {
			return err
		}

		if err := tx.Exec(`INSERT INTO pagerduty_credentials (created_at, updated_at, service_key,
               team_name, entity) VALUES(now(), now(), @service_key, @team_name, @entity)
               ON CONFLICT(team_name) 
			DO UPDATE SET "updated_at" = now(), service_key = excluded.service_key, entity = excluded.entity`,
			sql.Named("entity", credential.Entity),
			sql.Named("service_key", credential.PagerdutyCredentials),
			sql.Named("team_name", credential.TeamName),
		).Error; err != nil {
			return err
		}

		var credentials []SlackCredential
		tx.Model(SlackCredential{}).Where("entity= ?", credential.Entity).Find(&credentials)
		rows, err := tx.Raw(`select sw.team_name as team_name, pg.service_key as pg_service_key, sw.channel_name  as warning_channel_name, sw.webhook 						  as warning_webhook, sw.username as warning_username, 
 sc.channel_name as critical_channel_name, sc.webhook as critical_webhook, sc.username as critical_username
from slack_credentials as sw
 join slack_credentials as sc 
 on sc.team_name=sw.team_name 
 join pagerduty_credentials as pg
  on sc.team_name = pg.team_name
 where sw.entity= ? 
 and 
 sw.level='WARNING'
  and sc.entity=?
  and sc.level ='CRITICAL'  
 and pg.entity=?`, credential.Entity, credential.Entity, credential.Entity).Rows()
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
				&teamCredential.Slackcredentials.Warning.Channel, &teamCredential.Slackcredentials.Warning.Webhook, &teamCredential.Slackcredentials.Warning.Username,
				&teamCredential.Slackcredentials.Critical.Channel, &teamCredential.Slackcredentials.Critical.Webhook, &teamCredential.Slackcredentials.Critical.Username)
			teamCredentials[teamCredential.Name] = teamCredential

		}
		err = s.amClient.SyncConfig(alertmanager.EntityCredentials{
			Entity: credential.Entity,
			Teams:  teamCredentials,
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
				Channel:  warningSlackCredential.ChannelName,
				Webhook:  warningSlackCredential.Webhook,
				Username: warningSlackCredential.Username,
			},
			Critical: domain.SlackCredential{
				Channel:  criticalSlackCredential.ChannelName,
				Webhook:  criticalSlackCredential.Webhook,
				Username: criticalSlackCredential.Username,
			},
		},
	}
	return credential, nil
}

func NewService(db *gorm.DB, amClient alertmanager.Client) domain.AlertmanagerService {
	return &Service{repository: NewRepository(db),
		amClient: amClient,
	}
}
