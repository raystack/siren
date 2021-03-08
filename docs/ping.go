package docs

//-------------------------
// swagger:route GET /ping ping ping
// Ping call
// responses:
//   200: pingResponse

// Response body for Ping.
// swagger:response pingResponse
type pingResponse struct {
	// in:body
	Body string
}
