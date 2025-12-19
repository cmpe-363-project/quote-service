package restapi

import (
	"net/http"
	"quote-service/internal/repository"
	"quote-service/internal/restapi/routes"
	"quote-service/pkg/authorclient"
	"quote-service/pkg/logger"
	"strconv"
)

type App struct {
	Version      string
	Logger       logger.Logger
	Repository   repository.Repository
	AuthorClient *authorclient.Client

	Port int
	Host string
}

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (a *App) SetupAndRun() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/quote/{id}", routes.HandleGetQuoteByID(a.Logger, a.Repository, a.AuthorClient))
	mux.HandleFunc("GET /api/quote/random", routes.HandleGetRandomQuote(a.Logger, a.Repository, a.AuthorClient))
	mux.HandleFunc("GET /api/version", routes.HandleGetVersion(a.Version, a.AuthorClient, a.Logger))
	mux.HandleFunc("GET /api/mock-memory", routes.HandleAutoScalingDemo(a.Logger))

	// Wrap the mux with CORS middleware
	handler := corsMiddleware(mux)

	server := &http.Server{
		Addr:    a.Host + ":" + strconv.Itoa(a.Port),
		Handler: handler,
	}

	a.Logger.Info("Starting server", "host", a.Host, "port", strconv.Itoa(a.Port))
	if err := server.ListenAndServe(); err != nil {
		a.Logger.Error("Server failed to start")
		panic(err)
	}
}
