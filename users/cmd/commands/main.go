package main

import (
	"context"
	"os"

	"github.com/PratikKumar125/go-microservices/pkg/logging"
	"github.com/PratikKumar125/go-microservices/pkg/migrations/cli"
	"github.com/PratikKumar125/go-microservices/users/internal/config"
	"github.com/PratikKumar125/go-microservices/users/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
)

func main() {
	logger := logging.NewLogger("users", "info")

	var k = koanf.New(".")
	appConfig, err := config.NewAppConfig(k, "/Users/pratikkumar/Downloads/Personal/golang/Microservices/users/dev.env.yaml")
	if err != nil {
		panic(err)
	}

	dbInstance, err := db.NewDatabase(context.Background(), appConfig, logger)
	if err != nil {
		logger.Error("Failed to init db", "error", err)
		os.Exit(1)
	}
	defer dbInstance.Close()

	migrationsDir := "/Users/pratikkumar/Downloads/Personal/golang/Microservices/users/migrations"

	cli.RegisterDatabaseProvider(func() (*pgxpool.Pool, error) {
		return dbInstance.Pool(), nil
	})

	cli.InitRoot(&cli.Config{
		MigrationsDir: migrationsDir,
		Logger:        logger,
	})

	if err := cli.Execute(); err != nil {
		logger.Error("Command failed", "error", err)
		os.Exit(1)
	}
}
