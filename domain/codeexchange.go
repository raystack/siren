package domain

type OAuthPayload struct {
	Code      string `json:"code"`
	Workspace string `json:"workspace"`
}

type OAuthExchangeResponse struct {
	OK bool `json:"ok"`
}

type CodeExchangeService interface {
	Exchange(payload OAuthPayload) (*OAuthExchangeResponse, error)
	Migrate() error
}
