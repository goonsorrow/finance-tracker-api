package models

import (
	"errors"
	"time"
)

type Wallet struct {
	DisplayId int       `db:"display_id"`
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Balance   int64     `db:"computed_balance" json:"balance"`
	Currency  string    `db:"currency" json:"currency"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateWalletInput struct {
	Name           string  `json:"name" binding:"required"`
	InitialBalance float64 `json:"balance" binding:"omitempty,min=0"`
	Currency       string  `json:"currency" binding:"required,len=3"`
}

type UpdateWalletInput struct {
	Name    *string  `json:"name"`
	Balance *float64 `json:"balance" binding:"omitempty,min=0"`
}

type UpdateWalletData struct {
	Name    *string `json:"name"`
	Balance *int64  `json:"balance" binding:"omitempty,min=0"`
}

func (w CreateWalletInput) Validate() error {
	if w.Name == "" {
		return errors.New("name is required")
	}
	if w.InitialBalance < 0 {
		return errors.New("balance cannot be negative")
	}
	if w.Currency == "" || len(w.Currency) != 3 {
		return errors.New("currency must be 3 characters (USD, EUR, RUB)")
	}
	return nil
}

func (w UpdateWalletInput) Validate() error {
	if w.Name == nil && w.Balance == nil {
		return errors.New("at least one field must be provided for update")
	}
	if w.Balance != nil && *w.Balance < 0 {
		return errors.New("balance cannot be negative")
	}
	return nil
}
