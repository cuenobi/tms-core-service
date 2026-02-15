package main

import (
	"fmt"
	"os"

	"tms-core-service/cmd"
)

// @title TMS Core Service API
// @version 1.0
// @description Transportation Management System Core Service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@tms.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
