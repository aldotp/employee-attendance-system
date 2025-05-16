package main

import (
	"context"
	"log"

	"github.com/aldotp/employee-attendance-system/cmd/consumer"
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

	// define consumer command
	consumerCmd := cobra.Command{
		Use:   "consumer",
		Short: "Consumer is a command to start consumer worker",
	}

	consumerGenerateReporting := cobra.Command{
		Use:   "generate_reporting",
		Short: "Consumer is a command to start generate_reporting consumer server",
		Run: func(cmd *cobra.Command, args []string) {
			consumer.RunGenerateReportingAttendance(ctx)
		},
	}

	rootCmd.AddCommand(
		&restCmd,
		&consumerCmd,
	)

	consumerCmd.AddCommand(
		&consumerGenerateReporting)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}

}
