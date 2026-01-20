package models

import "time"

type TransactionFilter struct {
	WalletID  int
	Type      string // "income", "expense" или "internal"
	Category  string // для фильтрации
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	Offset    int
	SortBy    string // "date", "amount"
	SortOrder string // "asc", "desc"
}

type WalletFilter struct {
	UserID    int
	Currency  string
	Limit     int
	Offset    int
	SortBy    string // "name", "balance", "created_at"
	SortOrder string // "asc", "desc"
}
