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
	row := r.db.Raw(selectQuery, teamName).Row()
	if row.Err() != nil {
		return map[string]SlackCredential{}, PagerdutyCredential{}, row.Err()
	}
	var entity string
	pagerdutyCredential := PagerdutyCredential{}
	warningSlackCredential := SlackCredential{}
	criticalSlackCredential := SlackCredential{}
	row.Scan(&entity, &pagerdutyCredential.ServiceKey,
		&warningSlackCredential.ChannelName, &criticalSlackCredential.ChannelName)
	pagerdutyCredential.Entity = entity
	warningSlackCredential.Entity = entity
	criticalSlackCredential.Entity = entity
	return map[string]SlackCredential{"WARNING": warningSlackCredential,
		"CRITICAL": criticalSlackCredential}, pagerdutyCredential, nil
}
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
