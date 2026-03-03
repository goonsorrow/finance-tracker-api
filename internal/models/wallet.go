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
	Balance   int64     `db:"balance" json:"balance" example:"15000"`
	Currency  string    `db:"currency" json:"currency" example:"USD"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateWalletInput struct {
	Name           string  `json:"name" binding:"required" example:"Salary Card"`
	InitialBalance float64 `json:"balance" binding:"required" example:"1500"`
	Currency       string  `json:"currency" binding:"required,len=3" example:"USD"`
}

type UpdateWalletInput struct {
	Name     *string `json:"name" example:"Updated Wallet Name"`
	Currency *string `json:"currency" binding:"omitempty,len=3" example:"USD"`
}

func (w CreateWalletInput) Validate() error {
	if w.Name == "" {
		return errors.New("name is required")
	}
	if w.Currency == "" || len(w.Currency) != 3 {
		return errors.New("currency must be 3 characters (USD, EUR, RUB)")
	}
	return nil
}

func (w UpdateWalletInput) Validate() error {
	if w.Name == nil && w.Currency == nil {
		return errors.New("at least one field must be provided for update")
	}
	if w.Name != nil && *w.Name == "" {
		return errors.New("name must not be empty")
	}
	if w.Currency != nil && *w.Currency == "" {
		return errors.New("name must not be empty")
	}
	return nil
}
