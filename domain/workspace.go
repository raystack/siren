package domain

import "time"

type Workspace struct {
	Id        uint64    `json:"id"`
	Urn       string    `json:"urn"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkspaceService interface {
	ListWorkspaces() ([]Workspace, error)
	CreateWorkspace(*Workspace) (*Workspace, error)
	GetWorkspace(uint64) (*Workspace, error)
	UpdateWorkspace(*Workspace) (*Workspace, error)
	DeleteWorkspace(uint64) error
	Migrate() error
}
