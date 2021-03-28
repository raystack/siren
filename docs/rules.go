package docs

import "github.com/odpf/siren/domain"

// swagger:response rulesResponse
type rulesResponse struct {
	// in:body
	Body domain.Rule
}

//-------------------------
//-------------------------
// swagger:route GET /rules rules listRulesRequest
// List Rules API: This API lists all the existing rules with given filers in query params
// responses:
//   200: listRulesResponse

// swagger:parameters listRulesRequest
type listRulesRequest struct {
	// List Rule Request
	// in:query
	Namespace string `json:"namespace"`
	Entity    string `json:"entity"`
	GroupName string `json:"group_name"`
	Status    string `json:"status"`
	Template  string `json:"template"`
}

// List rules response
// swagger:response listRulesResponse
type listRulesResponse struct {
	// in:body
	Body []domain.Rule
}

//-------------------------
// swagger:route PUT /rules rules createRuleRequest
// Upsert Rule API: This API helps in creating a new rule or update an existing one with unique combination of namespace, entity, group_name, template
// responses:
//   200: rulesResponse

// swagger:parameters createRuleRequest
type createRuleRequest struct {
	// Create rule request
	// in:body
	Body domain.Rule
}
