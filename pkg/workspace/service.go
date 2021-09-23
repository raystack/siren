package workspace

import (
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository WorkspaceRepository
}

// NewService returns repository struct
func NewService(db *gorm.DB) domain.WorkspaceService {
	return &Service{NewRepository(db)}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) ListWorkspaces() ([]domain.Workspace, error) {
	workspaces, err := service.repository.List()
	if err != nil {
		return nil, err
	}

	domainWorkspaces := make([]domain.Workspace, 0, len(workspaces))
	for i := 0; i < len(workspaces); i++ {
		workspace, _ := workspaces[i].toDomain()
		domainWorkspaces = append(domainWorkspaces, *workspace)
	}
	return domainWorkspaces, nil

}

func (service Service) CreateWorkspace(workspace *domain.Workspace) (*domain.Workspace, error) {
	w := &Workspace{}
	w, err := w.fromDomain(workspace)
	if err != nil {
		return nil, err
	}

	newWorkspace, err := service.repository.Create(w)
	if err != nil {
		return nil, err
	}
	return newWorkspace.toDomain()
}

func (service Service) GetWorkspace(id uint64) (*domain.Workspace, error) {
	workspace, err := service.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return workspace.toDomain()
}

func (service Service) UpdateWorkspace(workspace *domain.Workspace) (*domain.Workspace, error) {
	w := &Workspace{}
	w, err := w.fromDomain(workspace)
	if err != nil {
		return nil, err
	}

	newWorkspace, err := service.repository.Update(w)
	if err != nil {
		return nil, err
	}
	return newWorkspace.toDomain()
}

func (service Service) DeleteWorkspace(id uint64) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
