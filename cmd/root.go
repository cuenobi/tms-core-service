package cmd

import (
	"fmt"
	"os"

	"tms-core-service/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	appCfg  *config.AppConfig
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "tms-core-service",
	Short: "TMS Core Service - Transportation Management System",
	Long: `TMS Core Service is a microservice built with Clean Architecture
for managing transportation and logistics operations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "env.yaml", "config file path")
}

// initConfig reads in config file
func initConfig() {
	var err error
	appCfg, err = config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
}

// GetConfig returns the loaded application configuration
func GetConfig() *config.AppConfig {
	return appCfg
}
