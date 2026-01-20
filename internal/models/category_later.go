package models

import (
	"errors"
	"time"
)

type Category struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name" binding:"required"`
	Type      string    `db:"type" json:"type"` // "initial", "income" или "expense"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Input для создания категории
type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=income expense"`
	Icon string `json:"icon"`
}

// Input для обновления категории
type UpdateCategoryInput struct {
	Name *string `json:"name"`
	Icon *string `json:"icon"`
}

func (c UpdateCategoryInput) Validate() error {
	if c.Name == nil && c.Icon == nil {
		return errors.New("at least one field must be provided for update")
	}
	return nil
}
