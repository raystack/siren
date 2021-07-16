package docs

import "github.com/odpf/siren/domain"

//-------------------------
//-------------------------
// swagger:route POST /code_exchange codeExchange postExchangeCodeRequest
// POST Code Exchange API
// This API exchanges oauth code with access token from slack server. client_id and client_secret are read from Siren ENV vars.
// responses:
//   200: postResponse

// swagger:parameters postExchangeCodeRequest
type postExchangeCodeRequest struct {
	// in:body
	Body domain.OAuthPayload
}

// POST codeExchange response
// swagger:response postResponse
type postResponse struct {
	// in:body
	Body domain.OAuthExchangeResponse
}
