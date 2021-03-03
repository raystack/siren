package templates

import (
	"fmt"
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Upsert(template *Template) (*Template, error) {
	var newTemplate, exisitingTemplate Template
	r.db.AutoMigrate(&Template{})
	//CREATE INDEX idx_tags on "templates" USING GIN ("tags");
	//SET enable_seqscan TO off;
	selectQuery := fmt.Sprintf(`select * from templates where name='%s';`, template.Name)
	result := r.db.Raw(selectQuery).Scan(&exisitingTemplate)
	if result.RowsAffected == 0 {
		result = r.db.Create(template)
	} else {
		result = r.db.Model(template).Where("id = ?", exisitingTemplate.ID).Updates(template)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	result = r.db.Raw(selectQuery).Scan(&newTemplate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newTemplate, nil
}
