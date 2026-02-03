package models

import (
	"errors"
	"time"
)

type Wallet struct {
	DisplayId int       `db:"display_id" json:"display_id"`
	ID        int       `db:"id" json:"id" example:"1"`
	UserID    int       `db:"user_id" json:"user_id" example:"10"`
	Name      string    `db:"name" json:"name" example:"Main Wallet"`
	Balance   int64     `db:"computed_balance" json:"balance" example:"15000"`
	Currency  string    `db:"currency" json:"currency" example:"USD"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateWalletInput struct {
	Name           string  `json:"name" binding:"required" example:"Salary Card"`
	InitialBalance float64 `json:"balance" binding:"omitempty,min=0" example:"1000.00"`
	Currency       string  `json:"currency" binding:"required,len=3" example:"USD"`
}

type UpdateWalletInput struct {
	Name    *string  `json:"name" example:"Updated Wallet Name"`
	Balance *float64 `json:"balance" binding:"omitempty,min=0" example:"500.50"`
}

type UpdateWalletData struct {
	Name    *string  `json:"name" example:"Updated Wallet Name"`
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
