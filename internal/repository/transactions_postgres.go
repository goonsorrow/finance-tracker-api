package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/goonsorrow/finance-tracker/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createTrQuery = `INSERT 
						INTO transactions (wallet_id, user_id, type, amount, category, description, date, created_at, updated_at) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) 
						RETURNING id`

	getAllTrQuery = `SELECT id, wallet_id, user_id, type, amount, category, description, date, created_at, updated_at  
						FROM transactions
						WHERE user_id = $1 AND wallet_id = $2
						ORDER BY date`

	getTrByIdQuery = `SELECT id, wallet_id, user_id, type, amount, category, description, date, created_at, updated_at  
						 FROM transactions
						 WHERE user_id = $1 AND wallet_id = $2 AND id = $3`

	deleteTrByIdQuery = `DELETE 
							FROM transactions 
        					WHERE user_id = $1	AND wallet_id = $2 AND id = $3`

	updateTrByIdQuery = `UPDATE transactions 
							SET type = COALESCE($1,type),
							amount = COALESCE($2,amount),
							category = COALESCE($3,category),
							description = COALESCE($4,description),
							date = COALESCE($5,date),
							updated_at = NOW()
							WHERE user_id = $6 AND wallet_id = $7 AND id = $8`
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) Create(ctx context.Context, userId, walletId int, input models.Transaction) (int, error) {
	var trId int

	row := r.db.QueryRowContext(ctx, createTrQuery,
		walletId,          //$1
		userId,            //$2
		input.Type,        //$3
		input.Amount,      //$4
		input.Category,    //$5
		input.Description, //$6
		input.Date)        //$7
	if err := row.Scan(&trId); err != nil {
		return 0, fmt.Errorf("failed to write down transaction: %w", err)
	}

	return trId, nil
}

func (r *TransactionPostgres) GetAll(ctx context.Context, userId, walletId int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.SelectContext(ctx, &transactions, getAllTrQuery,
		userId,   //$1
		walletId) //$2
	if err != nil {
		return nil, fmt.Errorf("failed gettin all transactions: %w", err)
	}
	return transactions, nil
}

func (r *TransactionPostgres) GetById(ctx context.Context, user_id, walletId, transactionId int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.GetContext(ctx, &transaction, getTrByIdQuery,
		user_id,       //$1
		walletId,      //$2
		transactionId) //$3
	if err != nil {
		return models.Transaction{}, fmt.Errorf("failed getting transaction: %w", err)
	}
	return transaction, nil
}

func (r *TransactionPostgres) Delete(ctx context.Context, userId, walletId, transactionId int) error {

	res, err := r.db.ExecContext(ctx, deleteTrByIdQuery,
		userId,        //$1
		walletId,      //$2
		transactionId) //$3

	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}
func (r *TransactionPostgres) Update(ctx context.Context, userId, walletId, transactionId int, input models.UpdateTransactionData) error {
	_, err := r.db.ExecContext(ctx, updateTrByIdQuery,
		input.Type,        // $1
		input.Amount,      // $2
		input.Category,    // $3
		input.Description, // $4
		input.Date,        // $5
		userId,            // $6
		walletId,          // $7
		transactionId)     // $8

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}
