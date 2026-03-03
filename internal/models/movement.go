package models

import (
	"errors"
	"time"
)

type Movement struct {
	ID          int       `db:"id" json:"id"`
	WalletID    int       `db:"wallet_id" json:"wallet_id"`
	UserId      int       `db:"user_id" json:"user_id"`
	Type        string    `db:"type" json:"type"` // "income" или "expense" или "initial"(только при создании кошелька с первоначальным балансом)
	Amount      int64     `db:"amount" json:"amount"`
	CategoryID  *int      `db:"category_id" json:"category_id"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Input для создания записи
type CreateMovementInput struct {
	Type        string    `json:"type" binding:"required,oneof=income expense initial" example:"expense"`
	Amount      float64   `json:"amount" binding:"required,gt=0" example:"150.50"`
	CategoryID  *int      `json:"category" binding:"required" example:"1"`
	Description string    `json:"description" example:"Grocery shopping"`
	Date        time.Time `json:"date" binding:"required" example:"2026-01-27T12:00:00Z"`
}

// Input для обновления операции
type UpdateMovementInput struct {
	Type        *string    `json:"type" example:"income"`
	Amount      *float64   `json:"amount" binding:"omitempty,gt=0" example:"200.00"`
	CategoryID  *int       `json:"category_id" example:"3"`
	Description *string    `json:"description" example:"Updated description"`
	Date        *time.Time `json:"date" example:"2026-01-28T15:00:00Z"`
}

type UpdateMovementData struct {
	Type        *string    `json:"type"`
	Amount      *int64     `json:"amount"`
	CategoryID  *int       `json:"category_id"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}

// Валидация
func (m CreateMovementInput) Validate() error {
	if m.Type != "income" && m.Type != "expense" && m.Type != "initial" {
		return errors.New("type must be 'income' or 'expense'")
	}
	if m.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if m.CategoryID == nil {
		return errors.New("category is required")
	}
	return nil
}

func (m UpdateMovementInput) Validate() error {
	if m.Type != nil && *m.Type != "income" && *m.Type != "expense" {
		return errors.New("type must be 'income' or 'expense'")
	}

	if m.Amount != nil && *m.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}
