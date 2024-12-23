package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/Justin-Fernbaugh/FinSight/handler"
)

const (
	serviceName = "FinSight"
)

var (
	command = &cobra.Command {
		Use: serviceName,
		Run: run,
	}
	clientID string
	token string
	databaseName string
)

func init() {
		command.Flags().StringVar(&clientID, "ynab-client-id", "", "The YNAB application client ID (required)")
		command.Flags().StringVar(&token, "ynab-token", "", "The YNAB client secret (required)")
		command.Flags().StringVar(&databaseName, "database-name", "", "The GCP Firestore database name (required)")

		// Mark the flags as required
		for _, flag := range []string{"ynab-client-id", "ynab-token", "database-name"} {
			err := command.MarkFlagRequired(flag)
			if err != nil {
				log.Fatalf("Error marking flag %s as required: %v", flag, err)
			}
		}
	
		// Bind flags to viper
		viper.BindPFlags(command.Flags())
	
		// Automatically read from environment variables
		viper.AutomaticEnv()
	
		// Replace dashes in flags with underscores for environment variable mapping
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func run(cmd *cobra.Command, args []string) {
	log.Println("Starting FinSight ...")

	if err := handler.NewHandler(token); err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}
}

func main() {
	// run the command
	command.Execute()
}