package models

import (
	"errors"
	"time"
)

type Category struct {
	ID         int       `db:"id" json:"id" example:"1"`
	UserID     *int      `db:"user_id" json:"user_id" binding:"omitempty" example:"10"` // Fixed typo: omitempy -> omitemty
	Name       string    `db:"name" json:"name" binding:"required" example:"Groceries"`
	Type       string    `db:"type" json:"type" example:"expense"` // "income" or "expense"
	Icon       *string   `db:"icon" json:"icon" binding:"omitempty" example:"🛒"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	UsageCount int       `db:"usage_count" json:"usage_count" example:"5"`
}

type CreateCategoryInput struct {
	Name string  `json:"name" binding:"required" example:"Salary"`
	Type string  `json:"type" binding:"required,oneof=income expense" example:"income"`
	Icon *string `json:"icon" example:"💰"`
}

type UpdateCategoryInput struct {
	Name *string `json:"name" binding:"omitempty" example:"Updated Category Name"`
	Icon *string `json:"icon" binding:"omitempty" example:"📝"`
}

func (c UpdateCategoryInput) Validate() error {
	if c.Name == nil && c.Icon == nil {
		return errors.New("at least one field must be provided for update")
	}
	return nil
}
