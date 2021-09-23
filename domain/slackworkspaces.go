package domain

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackWorkspaceService interface {
	GetChannels(string) ([]Channel, error)
}
