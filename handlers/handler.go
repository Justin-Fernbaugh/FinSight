package handlers

import (
	"context"
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/vertexai/genai"
	llm "github.com/Justin-Fernbaugh/FinSight/pkg"
	"github.com/brunomvsouza/ynab.go"
	ynabAPI "github.com/brunomvsouza/ynab.go/api"
	ynabBudget "github.com/brunomvsouza/ynab.go/api/budget"
	ynabTransaction "github.com/brunomvsouza/ynab.go/api/transaction"
)

func NewHandler(ynabClient ynab.ClientServicer, gemini *genai.GenerativeModel, tgUserID int64, daysBack int) error {
	transactions, err := getTransactions(ynabClient, daysBack)
	if err != nil {
		log.Fatalf("Error retrieving transactions: %v", err)
	}

	var unstructuredTransactions string
	for _, transaction := range transactions {
		unstructuredTransactions += createUnstructuredTransaction(transaction)
	}

	summary, err := llm.GenerateResponse(context.Background(), gemini, unstructuredTransactions)
	if err != nil {
		log.Fatalf("Error generating response: %v", err)
	}

	log.Printf("Response: %s", summary)
	msg := Message{
		UserID: tgUserID,
		Msg: summary,
	}
	SendMessagesAndShutdown([]Message{msg})
	
	return nil
}

func getTransactions(client ynab.ClientServicer, daysBack int) ([]*ynabTransaction.Transaction, error) {
	budgets, err := getBudgets(client)
	if err != nil {
		log.Fatalf("Error getting budgets: %v", err)
		return nil, err
	}
	budgets = budgets[:1] // Only use the first budget for now

	var transactions []*ynabTransaction.Transaction
	for _, budget := range budgets {
		budgetTransactions, err := getTransactionsByBudget(client, budget, daysBack)
		if err != nil {
			log.Fatalf("Error getting transactions: %v", err)
			return nil, err
		}
		transactions = append(transactions, budgetTransactions...)
	}
	return transactions, nil
}

// Simply return the first budget for now, later all could be aggregated or identified by something in the name.
func getBudgets(client ynab.ClientServicer) ([]*ynabBudget.Summary, error) {
	budgets, err := client.Budget().GetBudgets()
	if err != nil {
		log.Fatalf("Error getting budgets: %v", err)
		return nil, err
	}
	return budgets, nil
}

func getTransactionsByBudget(client ynab.ClientServicer, budget *ynabBudget.Summary, daysBack int) ([]*ynabTransaction.Transaction, error) {
	dateStr := time.Now().AddDate(0, 0, -daysBack).Format("2006-01-02")
	date, err := ynabAPI.DateFromString(dateStr)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
		return nil, err
	}

	transactions, err := client.Transaction().GetTransactions(budget.ID, &ynabTransaction.Filter{
		Since: &date,
	})
	if err != nil {
		log.Fatalf("Error getting transactions: %v", err)
	}
	return transactions, nil
}

func createUnstructuredTransaction(transaction *ynabTransaction.Transaction) string {
	var str string
	str += "date: " + transaction.Date.String() + "\n"
	str += "amount: " + convertAmount(transaction.Amount) + "\n"
	str += "account_name: " + transaction.AccountName + "\n"
	str += "payee_name: " + *transaction.PayeeName + "\n"
	str += "category_name: " + *transaction.CategoryName + "\n"
	str += "\n"
	return str
}

func convertAmount(amount int64) string {
	isNegative := amount < 0
	if isNegative {
		amount = -amount
	}

	amountStr := strconv.FormatInt(amount, 10)

	// Handle cases where the amount is less than 100 (i.e., less than $1.00)
	for len(amountStr) < 3 {
		amountStr = "0" + amountStr
	}

	dollars := amountStr[:len(amountStr)-3]
	cents := amountStr[len(amountStr)-2:]

	result := ""
	if isNegative {
		result = "-"
	}
	result += dollars + "." + cents

	return result
}