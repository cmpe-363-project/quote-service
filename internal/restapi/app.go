package restapi

import (
	"quote-service/internal/repository"
	"quote-service/internal/restapi/routes"
	"quote-service/pkg/authorclient"
	"quote-service/pkg/logger"
	"net/http"
	"strconv"
)

type App struct {
	Version      string
	Logger       logger.Logger
	Repository   repository.MysqlRepository
	AuthorClient *authorclient.Client

	Port int
	Host string
}

func (a *App) SetupAndRun() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/quote/{id}", routes.HandleGetQuoteByID(a.Logger, a.Repository, a.AuthorClient))
	mux.HandleFunc("GET /api/quote/random", routes.HandleGetRandomQuote(a.Logger, a.Repository, a.AuthorClient))
	mux.HandleFunc("GET /api/version", routes.HandleGetVersion(a.Version, a.AuthorClient, a.Logger))

	server := &http.Server{
		Addr:    a.Host + ":" + strconv.Itoa(a.Port),
		Handler: mux,
	}

	a.Logger.Info("Starting server", "host", a.Host, "port", strconv.Itoa(a.Port))
	if err := server.ListenAndServe(); err != nil {
		a.Logger.Error("Server failed to start")
		panic(err)
	}
}
