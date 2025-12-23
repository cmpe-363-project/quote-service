package main

import (
	hardcodedrepository "quote-service/internal/repository/hardcoded_adapter"
	"quote-service/internal/restapi"
	"quote-service/pkg/authorclient"
	"quote-service/pkg/logger/slog"
	"time"

	"github.com/caarlos0/env/v11"
)

type EnvVars struct {
	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`

	AuthorServiceURL string `env:"AUTHOR_SERVICE_URL,required"`
}

func main() {
	logger := slog.NewLogger(slog.NewLoggerArgs{
		LogFormat: "json",
	})

	envVars, err := env.ParseAs[EnvVars]()
	if err != nil {
		panic(err)
	}

	repo := hardcodedrepository.NewHardcodedRepository()

	authorClient := authorclient.NewClient(authorclient.NewClientConfig{
		BaseURL: envVars.AuthorServiceURL,
		Timeout: time.Second * 10,
	})

	app := &restapi.App{
		Version:      "v0.0.4",
		Logger:       logger,
		Repository:   repo,
		AuthorClient: authorClient,
		Port:         envVars.Port,
		Host:         envVars.Host,
	}
	app.SetupAndRun()
}
