package handlers

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/brunomvsouza/ynab.go"
	ynabAPI "github.com/brunomvsouza/ynab.go/api"
	ynabTransaction "github.com/brunomvsouza/ynab.go/api/transaction"
)

const (
	daysBack = 7
	budgetIdentifier = "Justin"
	tgUserID = 7053360498
)

func NewHandler(ynabClient ynab.ClientServicer) error {
	transactions, err := retrieveTransactions(ynabClient)
	if err != nil {
		log.Fatalf("Error retrieving transactions: %v", err)
	}

	var unstructuredTransactions string
	for _, transaction := range transactions {
		unstructuredTransactions += createUnstructuredTransaction(transaction)
	}
	msg := Message{
		UserID: tgUserID,
		Msg: unstructuredTransactions,
	}
	SendMessagesAndShutdown([]Message{msg})
	
	return nil
}

func retrieveTransactions(client ynab.ClientServicer) ([]*ynabTransaction.Transaction, error) {
	budgets, err := client.Budget().GetBudgets()
	if err != nil {
		log.Fatalf("Error getting budgets: %v", err)
		return nil, err
	}

	dateStr := time.Now().AddDate(0, 0, -daysBack).Format("2006-01-02")
	date, err := ynabAPI.DateFromString(dateStr)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
		return nil, err
	}

	allTransactions := make([]*ynabTransaction.Transaction, 0)
	for _, budget := range budgets {
		if !strings.Contains(budget.Name, budgetIdentifier) {
			continue
		}

		transactions, err := client.Transaction().GetTransactions(budget.ID, &ynabTransaction.Filter{
			Since: &date,
		})
		if err != nil {
			log.Fatalf("Error getting transactions: %v", err)
		}
		allTransactions = append(allTransactions, transactions...)
	}
	return allTransactions, nil
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