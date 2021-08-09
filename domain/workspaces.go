package domain

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WorkspaceService interface {
	GetChannels(string) ([]Channel, error)
}
