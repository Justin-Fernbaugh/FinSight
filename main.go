package main

import (
	"log"
	"strings"

	"github.com/Justin-Fernbaugh/FinSight/handlers"
	llm "github.com/Justin-Fernbaugh/FinSight/pkg"
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
	projectID string
	location string
	modelName string
	ynabClientID string
	ynabToken string
	tgBotToken string
	databaseName string
)

func init() {
		command.Flags().StringVar(&projectID, "project-id", "", "The GCP project ID (required)")
		command.Flags().StringVar(&ynabClientID, "ynab-client-id", "", "The YNAB application client ID (required)")
		command.Flags().StringVar(&ynabToken, "ynab-token", "", "The YNAB client secret (required)")
		command.Flags().StringVar(&databaseName, "database-name", "", "The GCP Firestore database name (required)")
		command.Flags().StringVar(&tgBotToken, "tg-bot-token", "", "The Telegram bot token (required)")
		command.Flags().StringVar(&location, "location", "us-west1", "The GCP location")
		command.Flags().StringVar(&modelName, "model-name", "gemini-1.5-flash-001", "The LLM model name")

		// Mark the flags as required
		for _, flag := range []string{"ynab-token", "tg-bot-token", "project-id"} {
			err := command.MarkFlagRequired(flag)
			if err != nil {
				log.Fatalf("Error marking flag %s as required: %v", flag, err)
			}
		}
		viper.BindPFlags(command.Flags())
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func run(cmd *cobra.Command, args []string) {
	log.Println("Starting FinSight ...")

    _, err := handlers.NewBotHandler(tgBotToken)
	if err != nil {
		log.Fatalf("Error creating bot handler: %v", err)
	}

	gemini, err := llm.Handler(projectID, location, modelName)
	if err != nil {
		log.Fatalf("Error creating LLM handler: %v", err)
	}

	ynabClient := ynab.NewClient(ynabToken)
	if err := handlers.NewHandler(ynabClient, gemini); err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	log.Println("FinSight exiting ...")
}

func main() {
	// run the command
	command.Execute()
}