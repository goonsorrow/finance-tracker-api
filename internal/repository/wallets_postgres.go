package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goonsorrow/finance-tracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type WalletPostgres struct {
	db *sqlx.DB
}

func NewWalletPostgres(db *sqlx.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

const (
	getAllQuery = `
    SELECT id, user_id, name, currency, computed_balance, created_at, updated_at,
           ROW_NUMBER() OVER (ORDER BY created_at DESC) AS display_id
    FROM wallets
    WHERE user_id = $1
    ORDER BY created_at DESC`

	getByIdQuery = `
    SELECT
      id, user_id, name, currency,
      (
        SELECT SUM(
          CASE
            WHEN t.type = 'income' OR t.type = 'initial' THEN +t.amount
            WHEN t.type = 'expense' THEN -t.amount
            ELSE 0
          END
        ) AS computed_balance
        FROM transactions t
        WHERE t.wallet_id = w.id
      ) AS computed_balance
    FROM wallets w
    WHERE w.user_id = $1 AND w.id = $2`

	createQuery = `
    INSERT INTO wallets (user_id, name, currency, computed_balance, created_at, updated_at)
    VALUES ($1, $2, $3, $4, NOW(), NOW())
    RETURNING id`

	createInitialTrQuery = `
    INSERT INTO transactions (wallet_id, user_id, type, amount, category, description, date, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`

	updateQuery = `
    UPDATE wallets
    SET name = COALESCE($1, name),
        computed_balance = COALESCE($2, computed_balance),
        updated_at = NOW()
    WHERE id = $3 AND user_id = $4`

	deleteQuery = `
    DELETE FROM wallets
    WHERE id = $1 AND user_id = $2`
)

func (r *WalletPostgres) Create(ctx context.Context, userId int, wallet models.Wallet) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, createQuery, userId, wallet.Name, wallet.Currency, wallet.Balance).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("[WalletCreate] failed to create wallet: %w", err)
	}

	return id, nil
}

func (r *WalletPostgres) CreateWithInitial(ctx context.Context, userId int, wallet models.Wallet) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("[createWithInitial] begin tx: %w", err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRowContext(ctx, createQuery, userId, wallet.Name, wallet.Currency, wallet.Balance).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("[createWithInitial] failed to create wallet: %w", err)
	}

	_, err = tx.ExecContext(ctx, createInitialTrQuery,
		id,                                   //$1
		userId,                               //$2
		"initial",                            //$3
		wallet.Balance,                       //$4
		"",                                   //$5
		"initial wallet balance transaction", //$6
		time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to write down transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to complete transaction %w", err)
	}

	return id, nil
}

func (r *WalletPostgres) GetAll(ctx context.Context, userId int) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.SelectContext(ctx, &wallets, getAllQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}
	return wallets, nil
}

func (r *WalletPostgres) GetById(ctx context.Context, userId, walletId int) (models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.GetContext(ctx, &wallet, getByIdQuery, userId, walletId)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("wallet not found: %w", err)
	}
	return wallet, nil
}

func (r *WalletPostgres) Update(ctx context.Context, userId, walletId int, input models.UpdateWalletData) error {
	result, err := r.db.ExecContext(ctx, updateQuery, input.Name, input.Balance, walletId, userId)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("wallet not found")
	}
	return nil
}

func (r *WalletPostgres) Delete(ctx context.Context, userId, walletId int) error {
	result, err := r.db.ExecContext(ctx, deleteQuery, walletId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("wallet not found")
	}
	return nil
}
