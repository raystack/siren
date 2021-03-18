package docs

import "github.com/odpf/siren/domain"

// swagger:response templatesResponse
type templatesResponse struct {
	// in:body
	Body domain.Template
}

//-------------------------
//-------------------------
// swagger:route GET /templates templates listTemplatesRequest
// List Templates API: This API lists all the existing templates with given filers in query params
// responses:
//   200: listResponse

// swagger:parameters listTemplatesRequest
type listTemplatesRequest struct {
	// List Template Request
	// in:query
	Tag string `json:"tag"`
}

// List templates response
// swagger:response listResponse
type listResponse struct {
	// in:body
	Body []domain.Template
}

//-------------------------
// swagger:route PUT /templates templates createTemplateRequest
// Upsert Templates API: This API helps in creating or updating a template with unique name
//responses:
//   200: templatesResponse

// swagger:parameters createTemplateRequest
type createTemplateRequest struct {
	// Create template request
	// in:body
	Body domain.Template
}

//-------------------------

// swagger:route GET /templates/{name} templates getTemplatesRequest
// Get Template API: This API gets a template given the template name
//responses:
//   200: templatesResponse

// swagger:parameters getTemplatesRequest
type getTemplatesRequest struct {
	// Get Template Request
	// in:path
	Name string `json:"name"`
}

//-------------------------

// swagger:route DELETE /templates/{name} templates deleteTemplatesRequest
// Delete Template API: This API deletes a template given the template name
// responses:
//   200: templatesResponse

// swagger:parameters deleteTemplatesRequest
type deleteTemplatesRequest struct {
	// Delete Template Request
	// in:path
	Name string `json:"name"`
}

//-------------------------

// swagger:route POST /templates/{name}/render templates renderTemplatesRequest
// Render Template API: This API renders the given template with given values
// responses:
//   200: renderTemplatesResponse

// swagger:parameters renderTemplatesRequest
type renderTemplatesRequest struct {
	// Render Template Request
	// in:path
	Name string `json:"name"`
}

// swagger:response renderTemplatesResponse
type renderTemplatesResponse struct {
	// Render Template Response
	// in:body
	Body string
}
