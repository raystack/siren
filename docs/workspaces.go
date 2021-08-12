package docs

import "github.com/odpf/siren/domain"

//-------------------------
// swagger:route GET /workspaces/{workspaceName}/channels workspaces getWorkspaceChannelsRequest
// Get Channels API: This API gets the list of joined channels within a slack workspace
// responses:
//   200: channelListResponse

// swagger:parameters getWorkspaceChannelsRequest
type getWorkspaceChannelsRequest struct {
	// name of the workspace
	// in:path
	WorkspaceName string `json:"workspaceName"`
}

// swagger:response channelListResponse
type channelListResponse struct {
	// in:body
	Body []domain.Channel
}
