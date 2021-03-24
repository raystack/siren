package alert

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r Repository) UpsertSlack(credential SlackCredential) error {
	//r.db.Clauses(clause.OnConflict{DoUpdates: true})
	result := r.db.Create(&credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r Repository) UpsertPagerduty(credential PagerdutyCredential) error {
	result := r.db.Create(&credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r Repository) GetCredential(teamName string) (map[string]SlackCredential, PagerdutyCredential, error) {
	row := r.db.Raw(`select sw.entity as entity, pg.service_key as pg_service_key, sw.channel_name  as warning_channel_name, sw.webhook  as warning_webhook, sw.username as warning_username, 
 sc.channel_name as critical_channel_name, sc.webhook as critical_webhook, sc.username as critical_username
from slack_credentials as sw
 join slack_credentials as sc 
 on sc.team_name=sw.team_name 
 join pagerduty_credentials as pg
  on sc.team_name = pg.team_name
 where sw.team_name= ? and 
 sw.level='WARNING'
  and sc.level ='CRITICAL'`, teamName).Row()
	if row.Err() != nil {
		return map[string]SlackCredential{}, PagerdutyCredential{}, row.Err()
	}
	var entity string
	pagerdutyCredential := PagerdutyCredential{}
	warningSlackCredential := SlackCredential{}
	criticalSlackCredential := SlackCredential{}
	row.Scan(&entity, &pagerdutyCredential.ServiceKey,
		&warningSlackCredential.ChannelName, &warningSlackCredential.Webhook, &warningSlackCredential.Username,
		&criticalSlackCredential.ChannelName, &criticalSlackCredential.Webhook, &criticalSlackCredential.Username)
	pagerdutyCredential.Entity = entity
	warningSlackCredential.Entity = entity
	criticalSlackCredential.Entity = entity
	return map[string]SlackCredential{"WARNING": warningSlackCredential,
		"CRITICAL": criticalSlackCredential}, pagerdutyCredential, nil
}
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
