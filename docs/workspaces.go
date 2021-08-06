package docs

type WorkspaceChannelsResponse []string

//-------------------------
// swagger:route GET /workspaces/{workspaceName}/channels workspaces getWorkspaceChannelsRequest
// Get Channels API: This API gets the list of joined channels within a slack workspace
//responses:
//   200: WorkspaceChannelsResponse

// swagger:parameters getWorkspaceChannelsRequest
type getWorkspaceChannelsRequest struct {
	// name of the workspace
	// in:path
	WorkspaceName string `json:"workspaceName"`
}
