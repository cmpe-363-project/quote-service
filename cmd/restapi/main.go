package main

import (
	"quote-service/internal/repository"
	"quote-service/internal/restapi"
	"quote-service/pkg/authorclient"
	"quote-service/pkg/logger/slog"
	"time"

	"github.com/caarlos0/env/v11"
	gormLogger "gorm.io/gorm/logger"
)

type EnvVars struct {
	Version string `env:"VERSION,required"`

	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`

	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

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

	gormLogLevel := gormLogger.Info

	repo, err := repository.NewMysqlRepository(repository.NewMysqlRepositoryConfig{
		Host:            envVars.DBHost,
		Port:            envVars.DBPort,
		Username:        envVars.DBUser,
		Password:        envVars.DBPassword,
		Database:        envVars.DBName,
		MaxOpenConns:    10,
		MaxIdleConns:    3,
		ConnMaxLifetime: time.Minute * 3,
		LogLevel:        &gormLogLevel,
		AutoMigrate:     true,
	})
	if err != nil {
		logger.Error("Failed to create MySQL repository", "error", err.Error())
		panic(err)
	}

	authorClient := authorclient.NewClient(authorclient.NewClientConfig{
		BaseURL: envVars.AuthorServiceURL,
		Timeout: time.Second * 10,
	})

	app := &restapi.App{
		Version:      envVars.Version,
		Logger:       logger,
		Repository:   *repo,
		AuthorClient: authorClient,
		Port:         envVars.Port,
		Host:         envVars.Host,
	}
	app.SetupAndRun()
}
