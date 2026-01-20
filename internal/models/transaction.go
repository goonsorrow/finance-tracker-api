package models

import (
	"errors"
	"time"
)

type Transaction struct {
	ID          int       `db:"id" json:"id"`
	WalletID    int       `db:"wallet_id" json:"wallet_id"`
	UserId      int       `db:"user_id" json:"user_id"`
	Type        string    `db:"type" json:"type"` // "income" или "expense" или "initial"(только при создании кошелька с первоначальным балансом)
	Amount      int64     `db:"amount" json:"amount"`
	Category    string    `db:"category" json:"category"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Input для создания транзакции
type CreateTransactionInput struct {
	Type        string    `json:"type" binding:"required,oneof=income expense initial"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Category    string    `json:"category" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// Input для обновления транзакции
type UpdateTransactionInput struct {
	Type        *string    `json:"type"`
	Amount      *float64   `json:"amount"`
	Category    *string    `json:"category"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}

// Валидация
func (t CreateTransactionInput) Validate() error {
	if t.Type != "income" && t.Type != "expense" {
		return errors.New("type must be 'income' or 'expense'")
	}
	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if t.Category == "" {
		return errors.New("category is required")
	}
	return nil
}

func (t UpdateTransactionInput) Validate() error {
	if t.Type != nil && *t.Type != "income" && *t.Type != "expense" {
		return errors.New("type must be 'income' or 'expense'")
	}

	if t.Amount != nil && *t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

type UpdateTransactionData struct {
	Type        *string    `json:"type"`
	Amount      *int64     `json:"amount"`
	Category    *string    `json:"category"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}
