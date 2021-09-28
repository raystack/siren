package workspace

import (
	"github.com/odpf/siren/domain"
	"time"
)

type Workspace struct {
	Id        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Urn       string `gorm:"index:idx_urn_name,unique"`
}

func (workspace *Workspace) fromDomain(t *domain.Workspace) (*Workspace, error) {
	workspace.Id = t.Id
	workspace.Urn = t.Urn
	workspace.Name = t.Name
	workspace.CreatedAt = t.CreatedAt
	workspace.UpdatedAt = t.UpdatedAt
	return workspace, nil
}

func (workspace *Workspace) toDomain() (*domain.Workspace, error) {
	if workspace == nil {
		return nil, nil
	}
	return &domain.Workspace{
		Id:        workspace.Id,
		Name:      workspace.Name,
		Urn:       workspace.Urn,
		CreatedAt: workspace.CreatedAt,
		UpdatedAt: workspace.UpdatedAt,
	}, nil
}

type WorkspaceRepository interface {
	Migrate() error
	List() ([]Workspace, error)
	Create(*Workspace) (*Workspace, error)
	Get(uint64) (*Workspace, error)
	Update(*Workspace) (*Workspace, error)
	Delete(uint64) error
}
