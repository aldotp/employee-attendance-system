package main

import (
	"context"
	"log"

	"github.com/aldotp/employee-attendance-system/cmd/http"
	"github.com/aldotp/employee-attendance-system/internal/adapter/config"
	"github.com/spf13/cobra"

	_ "github.com/joho/godotenv/autoload"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your_token}" to authenticate

func main() {
	ctx := context.Background()
	config.LoadConfig()
	rootCmd := &cobra.Command{}

	// Restful API Command
	restCmd := cobra.Command{
		Use:   "rest",
		Short: "Rest is a command to start Restful Api server",
		Run: func(cmd *cobra.Command, args []string) {
			http.RunHTTPServer(ctx)
		},
	}

	rootCmd.AddCommand(
		&restCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}

}
