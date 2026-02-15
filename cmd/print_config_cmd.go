package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// printConfigCmd represents the print-config command
var printConfigCmd = &cobra.Command{
	Use:   "print-config",
	Short: "Print the current configuration",
	Long:  `Prints the loaded configuration in JSON format for debugging purposes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return printConfig()
	},
}

func init() {
	rootCmd.AddCommand(printConfigCmd)
}

func printConfig() error {
	cfg := GetConfig()

	// Mask sensitive values
	maskedCfg := *cfg
	maskedCfg.Database.Password = "****"
	maskedCfg.Redis.Password = "****"
	maskedCfg.JWT.Secret = "****"

	jsonData, err := json.MarshalIndent(maskedCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println("Current Configuration:")
	fmt.Println(string(jsonData))

	return nil
}
