package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// newMigrationCmd represents the new-migration command
var newMigrationCmd = &cobra.Command{
	Use:   "new-migration [name]",
	Short: "Create a new migration file",
	Long:  `Creates a new pair of migration files (up and down) with a timestamp prefix.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return createMigration(args[0])
	},
}

func init() {
	rootCmd.AddCommand(newMigrationCmd)
}

func createMigration(name string) error {
	cfg := GetConfig()

	// Ensure migration directory exists
	if err := os.MkdirAll(cfg.Migration.Dir, 0o755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	// Generate timestamp-based filename
	timestamp := time.Now().Format("20060102150405")
	baseFilename := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join(cfg.Migration.Dir, baseFilename+".up.sql")
	downFile := filepath.Join(cfg.Migration.Dir, baseFilename+".down.sql")

	// Create up migration file
	if err := os.WriteFile(upFile, []byte("-- Add migration SQL here\n"), 0o644); err != nil {
		return fmt.Errorf("failed to create up migration: %w", err)
	}

	// Create down migration file
	if err := os.WriteFile(downFile, []byte("-- Add rollback SQL here\n"), 0o644); err != nil {
		return fmt.Errorf("failed to create down migration: %w", err)
	}

	fmt.Printf("✅ Created migration files:\n")
	fmt.Printf("   %s\n", upFile)
	fmt.Printf("   %s\n", downFile)

	return nil
}
