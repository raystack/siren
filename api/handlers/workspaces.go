package handlers

import (
	"github.com/gorilla/mux"
	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
	"net/http"
)

func GetWorkspaceChannels(service domain.WorkspaceService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		workspace := params["workspaceName"]

		result, err := service.GetChannels(workspace)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, result)
	}
}
