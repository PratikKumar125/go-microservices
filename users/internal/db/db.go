package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/PratikKumar125/go-microservices/pkg/logging"
	"github.com/PratikKumar125/go-microservices/users/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

var (
	instance *Database
	once     sync.Once
	initErr  error
)

func NewDatabase(ctx context.Context, cfg *config.AppConfig, logger *logging.Logger) (*Database, error) {
	once.Do(func() {
		// Construct connection string from config
		connString := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.ConfigService.String("db.host"), cfg.ConfigService.Int64("db.port"), cfg.ConfigService.String("db.user"), cfg.ConfigService.String("db.password"), cfg.ConfigService.String("db.name"), cfg.ConfigService.String("db.ssl_mode"),
		)

		// Initialize connection pool
		pool, err := pgxpool.New(ctx, connString)
		if err != nil {
			logger.Error("failed to initialize database pool", "error", err)
			return
		}

		// Create Database instance
		instance = &Database{pool: pool}

		// Ensure migrations table exists
		if err := instance.createMigrationsTable(ctx); err != nil {
			logger.Error("Failed to create migration table", "error", err)
			instance.Close()
			instance = nil
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}
	if instance == nil {
		return nil, fmt.Errorf("database initialization failed")
	}

	return instance, nil
}

func (db *Database) createMigrationsTable(ctx context.Context) error {
	const migrationTableQuery = `
	CREATE TABLE IF NOT EXISTS migrations (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		batch_no INT NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		down_migration TEXT NOT NULL DEFAULT ''
	);`
	_, err := db.pool.Exec(ctx, migrationTableQuery)
	if err != nil {
		return fmt.Errorf("failed to execute migrations table query: %w", err)
	}
	return nil
}

func (db *Database) Ping(ctx context.Context) error {
	if db.pool == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.pool = nil
	}
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}
