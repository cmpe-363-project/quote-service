package routes

import (
	"quote-service/internal/repository"
	restapiutils "quote-service/internal/restapi/utils"
	"quote-service/pkg/authorclient"
	"quote-service/pkg/logger"
	"net/http"
	"strconv"
)

// HandleGetQuoteByID
// /api/quote/{id}
func HandleGetQuoteByID(logger logger.Logger, repo repository.MysqlRepository, authorClient *authorclient.Client) http.HandlerFunc {
	type AuthorInfo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type Response struct {
		ID      int        `json:"id"`
		Message string     `json:"message"`
		Author  AuthorInfo `json:"author"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid quote ID", http.StatusBadRequest)
			return
		}

		quote, err := repo.GetQuoteByID(id)
		if err != nil {
			if err == repository.ErrNotFound {
				http.Error(w, "Quote not found", http.StatusNotFound)
				return
			}
			logger.ErrorWithCtx(r.Context(), "GetQuoteByID query failed", "error", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		authors, err := authorClient.GetAuthorsByIDs([]int{quote.AuthorID})
		if err != nil {
			logger.ErrorWithCtx(r.Context(), "Failed to get author", "error", err.Error())
			http.Error(w, "Failed to get author information", http.StatusInternalServerError)
			return
		}

		if len(authors) == 0 {
			logger.ErrorWithCtx(r.Context(), "Author not found", "authorID", strconv.Itoa(quote.AuthorID))
			http.Error(w, "Author not found", http.StatusNotFound)
			return
		}

		resp := Response{
			ID:      quote.ID,
			Message: quote.Message,
			Author: AuthorInfo{
				ID:   authors[0].ID,
				Name: authors[0].Name,
			},
		}

		restapiutils.WriteJSONResponse(w, http.StatusOK, resp)
	}
}

// HandleGetRandomQuote
// /api/quote/random
func HandleGetRandomQuote(logger logger.Logger, repo repository.MysqlRepository, authorClient *authorclient.Client) http.HandlerFunc {
	type AuthorInfo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type Response struct {
		ID      int        `json:"id"`
		Message string     `json:"message"`
		Author  AuthorInfo `json:"author"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := repo.GetRandomQuote()
		if err != nil {
			if err == repository.ErrNotFound {
				http.Error(w, "No quotes found", http.StatusNotFound)
				return
			}
			logger.ErrorWithCtx(r.Context(), "GetRandomQuote query failed", "error", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		authors, err := authorClient.GetAuthorsByIDs([]int{quote.AuthorID})
		if err != nil {
			logger.ErrorWithCtx(r.Context(), "Failed to get author", "error", err.Error())
			http.Error(w, "Failed to get author information", http.StatusInternalServerError)
			return
		}

		if len(authors) == 0 {
			logger.ErrorWithCtx(r.Context(), "Author not found", "authorID", strconv.Itoa(quote.AuthorID))
			http.Error(w, "Author not found", http.StatusNotFound)
			return
		}

		resp := Response{
			ID:      quote.ID,
			Message: quote.Message,
			Author: AuthorInfo{
				ID:   authors[0].ID,
				Name: authors[0].Name,
			},
		}

		restapiutils.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
