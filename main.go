package main

import (
	"log"
	"strings"

	"github.com/Justin-Fernbaugh/FinSight/handlers"
	"github.com/brunomvsouza/ynab.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	serviceName = "FinSight"
)

var (
	command = &cobra.Command {
		Use: serviceName,
		Run: run,
	}
	ynabClientID string
	ynabToken string
	tgBotToken string
	databaseName string
)

func init() {
		command.Flags().StringVar(&ynabClientID, "ynab-client-id", "", "The YNAB application client ID (required)")
		command.Flags().StringVar(&ynabToken, "ynab-token", "", "The YNAB client secret (required)")
		command.Flags().StringVar(&databaseName, "database-name", "", "The GCP Firestore database name (required)")
		command.Flags().StringVar(&tgBotToken, "tg-bot-token", "", "The Telegram bot token (required)")

		// Mark the flags as required
		for _, flag := range []string{"ynab-client-id", "ynab-token", "tg-bot-token", "database-name"} {
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

    _, err := handlers.NewBotHandler(tgBotToken)
	if err != nil {
		log.Fatalf("Error creating bot handler: %v", err)
	}

	ynabClient := ynab.NewClient(ynabToken)
	if err := handlers.NewHandler(ynabClient); err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	log.Println("FinSight exiting ...")
}

func main() {
	// run the command
	command.Execute()
}