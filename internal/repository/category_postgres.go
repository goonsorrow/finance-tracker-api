package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

const (
	createCategoryQuery = `INSERT INTO categories (name, type, user_id) 
								VALUES ($1,$2,$3)
								RETURNING id`

	getAllCategoriesQuery = `SELECT * 
								FROM categories 
								WHERE user_id = $1 OR user_id IS NULL
								ORDER BY CASE WHEN user_id = $1 THEN 0 ELSE 1 END,
								usage_count DESC, name`

	getCategoryByIdQuery = `SELECT *
								FROM categories 
								WHERE id = $1 
								AND (user_id = $2 OR user_id IS NULL)`

	updateCategoryById = `UPDATE categories
							SET name = COALESCE($1,name),
								icon = COALESCE($2,icon)
							WHERE id = $3 AND (user_id = $4 OR user_id IS NULL)`

	deleteCategoryById = `DELETE
							FROM categories
							WHERE id = $1 AND user_id = $2`
)

func (r CategoryPostgres) Create(ctx context.Context, userId int, input models.Category) (int, error) {
	var id int
	row := r.db.QueryRowContext(ctx, createCategoryQuery,
		input.Name, //$1
		input.Type, //$2
		userId)     //$3
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed creating category:%w", err)
	}

	return id, nil
}

func (r CategoryPostgres) GetAll(ctx context.Context, userId int) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.SelectContext(ctx, &categories, getAllCategoriesQuery, userId) //$1
	if err != nil {
		return nil, fmt.Errorf("failed getting all user categories: %w", err)
	}
	return categories, nil
}

func (r CategoryPostgres) GetById(ctx context.Context, userId int, categoryId int) (models.Category, error) {
	var category models.Category
	err := r.db.GetContext(ctx, &category, getCategoryByIdQuery,
		categoryId, //$1
		userId)     //$2
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Category{}, ErrRecordNotFound
		}
		return models.Category{}, fmt.Errorf("failed getting category by id: %w", err)
	}

	return category, nil
}

func (r CategoryPostgres) Update(ctx context.Context, userId, categoryId int, input models.UpdateCategoryInput) error {
	res, err := r.db.ExecContext(ctx, updateCategoryById, input.Name, input.Icon, categoryId, userId)
	if err != nil {
		return fmt.Errorf("failed to update category:%w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r CategoryPostgres) Delete(ctx context.Context, userId, categoryId int) error {
	res, err := r.db.ExecContext(ctx, deleteCategoryById, categoryId, userId)
	if err != nil {
		return fmt.Errorf("failed deleting category:%w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}
