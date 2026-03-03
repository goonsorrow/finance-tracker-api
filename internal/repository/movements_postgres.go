package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createMQuery = `INSERT 
						INTO movements (wallet_id, user_id, type, amount, category_id, description, date, created_at, updated_at) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) 
						RETURNING id`

	getAllMQuery = `SELECT id, wallet_id, user_id, type, amount, category_id, description, date, created_at, updated_at  
						FROM movements
						WHERE user_id = $1 AND wallet_id = $2
						ORDER BY date`

	getMByIdQuery = `SELECT id, wallet_id, user_id, type, amount, category_id, description, date, created_at, updated_at  
						 FROM movements
						 WHERE user_id = $1 AND wallet_id = $2 AND id = $3`

	deleteMByIdQuery = `DELETE 
							FROM movements 
        					WHERE user_id = $1	AND wallet_id = $2 AND id = $3`

	updateMByIdQuery = `UPDATE movements 
							SET type = COALESCE($1,type),
							amount = COALESCE($2,amount),
							category_id = COALESCE($3,category_id),
							description = COALESCE($4,description),
							date = COALESCE($5,date),
							updated_at = NOW()
							WHERE user_id = $6 AND wallet_id = $7 AND id = $8`
)

type MovementPostgres struct {
	db         *sqlx.DB
	transactor Transactor
}

func NewMovementPostgres(db *sqlx.DB, transactor Transactor) *MovementPostgres {
	return &MovementPostgres{db: db, transactor: transactor}
}

func (r *MovementPostgres) Create(ctx context.Context, userId, walletId int, input models.Movement) (int, error) {
	var mId int

	exc := r.transactor.GetExecutor(ctx)

	err := exc.QueryRowxContext(ctx, createMQuery,
		walletId,              //$1
		userId,                //$2
		input.Type,            //$3
		input.Amount,          //$4
		input.CategoryID,      //$5
		input.Description,     //$6
		input.Date).Scan(&mId) //$7
	if err != nil {
		return 0, fmt.Errorf("[MovementPostgres.Create] failed to write down movement: %w", err)
	}

	return mId, nil
}

func (r *MovementPostgres) GetAll(ctx context.Context, userId, walletId int) ([]models.Movement, error) {
	var movements []models.Movement

	exc := r.transactor.GetExecutor(ctx)

	err := sqlx.SelectContext(ctx, exc, &movements, getAllMQuery,
		userId,   //$1
		walletId) //$2
	if err != nil {
		return nil, fmt.Errorf("failed gettin all movements: %w", err)
	}
	return movements, nil
}

func (r *MovementPostgres) GetById(ctx context.Context, user_id, walletId, movementId int) (models.Movement, error) {
	var movement models.Movement
	exc := r.transactor.GetExecutor(ctx)

	err := sqlx.GetContext(ctx, exc, &movement, getMByIdQuery,
		user_id,    //$1
		walletId,   //$2
		movementId) //$3
	if err != nil {
		return models.Movement{}, fmt.Errorf("failed getting movement: %w", err)
	}
	return movement, nil
}

func (r *MovementPostgres) Delete(ctx context.Context, userId, walletId, movementId int) error {
	exc := r.transactor.GetExecutor(ctx)

	res, err := exc.ExecContext(ctx, deleteMByIdQuery,
		userId,     //$1
		walletId,   //$2
		movementId) //$3

	if err != nil {
		return fmt.Errorf("[MovementPostgres.Delete] failed to delete movement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("[MovementPostgres.Delete] failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("[MovementPostgres.Delete] movement not found")
	}
	return nil
}
func (r *MovementPostgres) Update(ctx context.Context, userId, walletId, movementId int, input models.UpdateMovementData) error {
	exc := r.transactor.GetExecutor(ctx)

	_, err := exc.ExecContext(ctx, updateMByIdQuery,
		input.Type,        // $1
		input.Amount,      // $2
		input.CategoryID,  // $3
		input.Description, // $4
		input.Date,        // $5
		userId,            // $6
		walletId,          // $7
		movementId)        // $8
	if err != nil {
		return fmt.Errorf("[MovementPostgres.Update] failed to update movement: %w", err)
	}
	return nil
}
