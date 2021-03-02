package store

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
)

const (
	tableNameTemplates = "templates"
)

var (
	queryCreateTemplate = fmt.Sprintf(`
		INSERT INTO %s
		(name, body, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?) Returning *;
	`, tableNameTemplates)
)

type TemplatesStore struct {
	db *gorm.DB
}

func NewTemplatesStore(db *gorm.DB) *TemplatesStore {
	return &TemplatesStore{db: db}
}

type Template struct {
	gorm.Model
	Name string         `json:"name" gorm:"index:idx_name,unique"`
	Body string         `json:"body"`
	Tags pq.StringArray `gorm:"type:text[]" json:"tags"`
}

//UpsertTemplate upserts tempaltes based on template name
func (store *TemplatesStore) Upsert(template *Template) (*domain.Template, error) {
	var newTemplate, exisitingTemplate domain.Template
	store.db.AutoMigrate(&Template{})
	//CREATE INDEX idx_tags on "templates" USING GIN ("tags");
	//SET enable_seqscan TO off;
	selectQuery := fmt.Sprintf(`select * from templates where name='%s';`, template.Name)
	result := store.db.Raw(selectQuery).Scan(&exisitingTemplate)
	if result.RowsAffected == 0 {
		result = store.db.Create(template)
	} else {
		result = store.db.Model(template).Where("id = ?", exisitingTemplate.ID).Updates(template)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	result = store.db.Raw(selectQuery).Scan(&newTemplate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newTemplate, nil
}
