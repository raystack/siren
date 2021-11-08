package docs

import "github.com/odpf/siren/domain"

//-------------------------
//-------------------------
// swagger:route GET /history alertHistory getAlertHistoryRequest
// GET Alert History API: This API lists stored alert history for given filers in query params
// responses:
//   200: getResponse

// swagger:parameters getAlertHistoryRequest
type getAlertHistoryRequest struct {
	// in:query
	Resource  string `json:"resource"`
	StartTime uint32 `json:"startTime"`
	EndTime   uint32 `json:"endTime"`
}

// Get alertHistory response
// swagger:response getResponse
type getResponse struct {
	// in:body
	Body []domain.Alerts
}

//-------------------------
// swagger:route POST /history alertHistory createAlertHistoryRequest
// Create Alert History API: This API create alert history

// swagger:parameters createAlertHistoryRequest
type createAlertHistoryRequest struct {
	// in:body
	Body []domain.Alerts
}
