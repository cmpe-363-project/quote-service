package routes

import (
	"quote-service/pkg/authorclient"
	restapiutils "quote-service/internal/restapi/utils"
	"quote-service/pkg/logger"
	"net/http"
)

// HandleGetVersion
// /api/version
func HandleGetVersion(version string, authorClient *authorclient.Client, logger logger.Logger) http.HandlerFunc {
	type Response struct {
		QuoteService  string `json:"quote-service"`
		AuthorService string `json:"author-service"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authorVersion, err := authorClient.GetVersion()
		if err != nil {
			logger.ErrorWithCtx(r.Context(), "Failed to get author-service version", "error", err.Error())
			http.Error(w, "Failed to get author-service version", http.StatusInternalServerError)
			return
		}

		restapiutils.WriteJSONResponse(w, http.StatusOK, Response{
			QuoteService:  version,
			AuthorService: authorVersion,
		})
	}
}
