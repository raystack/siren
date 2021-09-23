package workspace

import (
	"errors"
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

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Workspace{})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) List() ([]Workspace, error) {
	var workspaces []Workspace
	selectQuery := fmt.Sprintf("select * from workspaces")
	result := r.db.Raw(selectQuery).Find(&workspaces)
	if result.Error != nil {
		return nil, result.Error
	}
	return workspaces, nil
}

func (r Repository) Create(workspace *Workspace) (*Workspace, error) {
	var newWorkspace Workspace

	result := r.db.Create(workspace)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Where(fmt.Sprintf("urn = '%s'", workspace.Urn)).Find(&newWorkspace)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newWorkspace, nil
}

func (r Repository) Get(id uint64) (*Workspace, error) {
	var workspace Workspace
	result := r.db.Where(fmt.Sprintf("id = '%d'", id)).Find(&workspace)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no workspace found")
	}

	return &workspace, nil
}

func (r Repository) Update(workspace *Workspace) (*Workspace, error) {
	var newWorkspace, existingWorkspace Workspace
	result := r.db.Where(fmt.Sprintf("id = '%d'", workspace.Id)).Find(&existingWorkspace)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("workspace doesn't exits")
	} else {
		result = r.db.Where("id = ?", workspace.Id).Updates(workspace)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Where(fmt.Sprintf("id = '%d'", workspace.Id)).Find(&newWorkspace)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newWorkspace, nil
}

func (r Repository) Delete(id uint64) error {
	var workspace Workspace
	result := r.db.Where("id = ?", id).Delete(&workspace)
	return result.Error
}
