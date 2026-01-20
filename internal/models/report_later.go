package models

import "time"

type SummaryReport struct {
	TotalBalance     float64 `json:"total_balance"`
	TotalIncome      float64 `json:"total_income"`
	TotalExpense     float64 `json:"total_expense"`
	NetBalance       float64 `json:"net_balance"` // income - expense
	WalletCount      int     `json:"wallet_count"`
	TransactionCount int     `json:"transaction_count"`
}

type CategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
	Percent  float64 `json:"percent"`
}

type CategoryReport struct {
	Data       []CategorySummary `json:"data"`
	GrandTotal float64           `json:"grand_total"`
}

type MonthlySummary struct {
	Month     string  `json:"month"`
	Income    float64 `json:"income"`
	Expense   float64 `json:"expense"`
	NetChange float64 `json:"net_change"`
}

type MonthlyReport struct {
	Data []MonthlySummary `json:"data"`
}

type WalletStats struct {
	WalletID         int       `json:"wallet_id"`
	WalletName       string    `json:"wallet_name"`
	Balance          float64   `json:"balance"`
	Currency         string    `json:"currency"`
	TotalIncome      float64   `json:"total_income"`
	TotalExpense     float64   `json:"total_expense"`
	TransactionCount int       `json:"transaction_count"`
	LastTransaction  time.Time `json:"last_transaction"`
}
