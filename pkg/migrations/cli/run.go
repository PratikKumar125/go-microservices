package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var databaseProvider func() (*pgxpool.Pool, error)

var runMigrationsCmd = &cobra.Command{
	Use:   "run:migrations",
	Short: "Run pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		db, err := getDatabase()
		if err != nil {
			cfg.Logger.Error("Failed to get database", "error", err)
			return
		}

		cfg.Logger.Info("Starting migrations")

		completed, err := getCompletedMigrations(ctx, db)
		if err != nil {
			cfg.Logger.Error("Failed to get completed migrations", "error", err)
			return
		}

		pending, err := findPendingMigrations(completed)
		if err != nil {
			cfg.Logger.Error("Failed to find pending migrations", "error", err)
			return
		}

		if len(pending) == 0 {
			cfg.Logger.Info("No pending migrations")
			return
		}

		if err := executeMigrations(ctx, db, pending); err != nil {
			cfg.Logger.Error("Failed to execute migrations", "error", err)
			return
		}
		cfg.Logger.Info("Applied migrations", "count", len(pending))
	},
}

func RegisterDatabaseProvider(provider func() (*pgxpool.Pool, error)) {
	databaseProvider = provider
}

func getDatabase() (*pgxpool.Pool, error) {
	return databaseProvider()
}

func getCompletedMigrations(ctx context.Context, db *pgxpool.Pool) (map[string]bool, error) {
	completed := make(map[string]bool)
	rows, err := db.Query(ctx, "SELECT name FROM migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		completed[name] = true
	}
	return completed, nil
}

func findPendingMigrations(completed map[string]bool) ([]string, error) {
	if _, err := os.Stat(cfg.MigrationsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations dir missing: %w", err)
	}

	var migrations []string
	err := filepath.WalkDir(cfg.MigrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".up.sql") {
			name := strings.TrimSuffix(d.Name(), ".up.sql")
			if !completed[name] {
				migrations = append(migrations, d.Name())
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations: %w", err)
	}

	sort.Strings(migrations)
	return migrations, nil
}

func executeMigrations(ctx context.Context, db *pgxpool.Pool, migrations []string) error {
	batchNo, err := getLatestBatchNo(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get batch number: %w", err)
	}
	batchNo++

	for _, m := range migrations {
		cfg.Logger.Info("Applying", "file", m)
		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin tx for %s: %w", m, err)
		}

		sqlBytes, err := os.ReadFile(filepath.Join(cfg.MigrationsDir, m))
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to read %s: %w", m, err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to execute %s: %w", m, err)
		}

		downFile := strings.TrimSuffix(m, ".up.sql") + ".down.sql"
		name := strings.TrimSuffix(m, ".up.sql")
		query := `INSERT INTO migrations (name, batch_no, down_migration) VALUES ($1, $2, $3)`
		if _, err := tx.Exec(ctx, query, name, batchNo, downFile); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record %s: %w", m, err)
		}

		if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to commit %s: %w", m, err)
		}
		cfg.Logger.Info("Applied", "file", m)
	}
	return nil
}

func getLatestBatchNo(ctx context.Context, db *pgxpool.Pool) (int, error) {
	var batchNo int
	rows, err := db.Query(ctx, "SELECT COALESCE(MAX(batch_no), 0) FROM migrations")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch batch: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&batchNo); err != nil {
			return 0, fmt.Errorf("failed to scan: %w", err)
		}
	}
	return batchNo, nil
}
