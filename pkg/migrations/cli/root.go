package cli

import (
	"os"
	"github.com/spf13/cobra"
)

var cfg *Config

var rootCmd = &cobra.Command{
	Use:   "gopher",
	Short: "Gopher CLI for migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(cfg.MigrationsDir, 0755); err != nil {
			cfg.Logger.Error("Failed to create migrations dir", "error", err)
			os.Exit(1)
		}
	},
}

// InitRoot initializes the root command.
func InitRoot(config *Config) {
	cfg = config
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(runMigrationsCmd)
	rootCmd.AddCommand(rollbackCmd)
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}