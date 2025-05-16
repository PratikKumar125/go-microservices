package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var steps int

var rollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback migrations",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		db, err := getDatabase()
		if err != nil {
			cfg.Logger.Error("Failed to get database", "error", err)
			return
		}

		cfg.Logger.Info("Starting rollback", "steps", steps)

		if err := rollbackMigrations(ctx, db, steps); err != nil {
			cfg.Logger.Error("Failed to rollback", "error", err)
			return
		}
		cfg.Logger.Info("Rollback done")
	},
}

func init() {
	rollbackCmd.Flags().IntVar(&steps, "steps", 1, "Batches to rollback")
}

func rollbackMigrations(ctx context.Context, db *pgxpool.Pool, steps int) error {
	latestBatch, err := getLatestBatchNo(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get batch: %w", err)
	}
	if latestBatch == 0 {
		cfg.Logger.Info("No migrations to rollback")
		return nil
	}

	targetBatch := latestBatch - steps
	if targetBatch < 0 {
		targetBatch = 0
	}

	type migration struct {
		name          string
		downMigration string
		batchNo       int
	}

	var toRollback []migration
	rows, err := db.Query(ctx, "SELECT name, down_migration, batch_no FROM migrations WHERE batch_no > $1 ORDER BY batch_no DESC, executed_at DESC", targetBatch)
	if err != nil {
		return fmt.Errorf("failed to fetch migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m migration
		if err := rows.Scan(&m.name, &m.downMigration, &m.batchNo); err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}
		toRollback = append(toRollback, m)
	}

	sort.Slice(toRollback, func(i, j int) bool {
		if toRollback[i].batchNo == toRollback[j].batchNo {
			return toRollback[i].name > toRollback[j].name
		}
		return toRollback[i].batchNo > toRollback[j].batchNo
	})

	for _, m := range toRollback {
		cfg.Logger.Info("Rolling back", "name", m.name)
		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin tx for %s: %w", m.name, err)
		}

		sqlBytes, err := os.ReadFile(filepath.Join(cfg.MigrationsDir, m.downMigration))
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to read %s: %w", m.downMigration, err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to execute %s: %w", m.downMigration, err)
		}

		if _, err := tx.Exec(ctx, "DELETE FROM migrations WHERE name = $1", m.name); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to delete %s: %w", m.name, err)
		}

		if err := tx.Commit(ctx); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to commit %s: %w", m.name, err)
		}
		cfg.Logger.Info("Rolled back", "name", m.name)
	}
	return nil
}
