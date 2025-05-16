package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/spf13/cobra"
)

var (
	migrationName string
)

var migrateCmd = &cobra.Command{
	Use:   "make:migration",
	Short: "Create migration files",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Logger.Info("Creating migration", "name", migrationName)

		if err := os.MkdirAll(cfg.MigrationsDir, 0755); err != nil {
			cfg.Logger.Error("Error creating migrations dir", "error", err)
			return
		}

		timestamp := time.Now().Format("20060102150405")
		upFile := filepath.Join(cfg.MigrationsDir, fmt.Sprintf("%s_%s.up.sql", timestamp, migrationName))
		downFile := filepath.Join(cfg.MigrationsDir, fmt.Sprintf("%s_%s.down.sql", timestamp, migrationName))

		if err := createFile(upFile); err != nil {
			cfg.Logger.Error("Error creating up file", "file", upFile, "error", err)
			return
		}
		cfg.Logger.Info("Created", "file", upFile)

		if err := createFile(downFile); err != nil {
			cfg.Logger.Error("Error creating down file", "file", downFile, "error", err)
			return
		}
		cfg.Logger.Info("Created", "file", downFile)
	},
}

func init() {
	migrateCmd.Flags().StringVar(&migrationName, "name", "", "Migration name")
	migrateCmd.MarkFlagRequired("name")
}

func createFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	return file.Close()
}