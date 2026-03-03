package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type WalletPostgres struct {
	db         *sqlx.DB
	transactor Transactor
}

func NewWalletPostgres(db *sqlx.DB, transactor Transactor) *WalletPostgres {
	return &WalletPostgres{db: db, transactor: transactor}
}

const (
	getAllQuery = `
        SELECT 
            id, user_id, name, currency, balance, created_at, updated_at
        FROM wallets
        WHERE user_id = $1
        ORDER BY created_at DESC`

	getByIdQuery = `
        SELECT 
            id, user_id, name, currency, balance, created_at, updated_at
        FROM wallets
        WHERE user_id = $1 AND id = $2`

	createQuery = `
        INSERT INTO wallets (user_id, name, currency, balance, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id`

	createInitialTrQuery = `
        INSERT INTO movements (wallet_id, user_id, type, amount, category_id, description, date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`

	addToBalanceQuery = `
        UPDATE wallets
        SET balance = balance + $1,
            updated_at = NOW()
        WHERE id = $2`

	updateQuery = `
        UPDATE wallets
        SET name = COALESCE($1, name),
			currency = COALESCE($2, currency),
            updated_at = NOW()
        WHERE id = $3 AND user_id = $4`

	deleteQuery = `
        DELETE FROM wallets
        WHERE id = $1 AND user_id = $2`
)

func (r *WalletPostgres) Create(ctx context.Context, userId int, wallet models.Wallet) (int, error) {
	var id int
	err := r.transactor.GetExecutor(ctx).QueryRowxContext(ctx, createQuery, userId, wallet.Name, wallet.Currency, wallet.Balance).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("[WalletPostgres.Create] failed to create empty wallet: %w", err)
	}
	return id, nil
}

func (r *WalletPostgres) GetAll(ctx context.Context, userId int) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := sqlx.SelectContext(ctx, r.transactor.GetExecutor(ctx), &wallets, getAllQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("[WalletPostgres.GetAll] failed to get wallets: %w", err)
	}
	return wallets, nil
}

func (r *WalletPostgres) GetById(ctx context.Context, userId, walletId int) (models.Wallet, error) {
	var wallet models.Wallet
	err := sqlx.GetContext(ctx, r.transactor.GetExecutor(ctx), &wallet, getByIdQuery, userId, walletId)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("[WalletPostgres.GetById] wallet not found: %w", err)
	}
	return wallet, nil
}

func (r *WalletPostgres) Update(ctx context.Context, userId, walletId int, input models.UpdateWalletInput) error {
	var name, currency interface{} = nil, nil
	if input.Name != nil {
		name = *input.Name
	}
	if input.Currency != nil {
		currency = *input.Currency
	}
	result, err := r.transactor.GetExecutor(ctx).ExecContext(ctx, updateQuery, name, currency, walletId, userId)
	if err != nil {
		return fmt.Errorf("[WalletPostgres.Update] failed to update wallet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[WalletPostgres.Update] failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("[WalletPostgres.Update] wallet not found")
	}
	return nil
}

func (r *WalletPostgres) AddToBalance(ctx context.Context, walletId int, deltaCents int64) error {
	result, err := r.transactor.GetExecutor(ctx).ExecContext(ctx, addToBalanceQuery, deltaCents, walletId)
	if err != nil {
		return fmt.Errorf("[WalletPostgres.AddToBalance] failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[WalletPostgres.AddToBalance] failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("[WalletPostgres.AddToBalance] wallet not found")
	}
	return nil
}

func (r *WalletPostgres) Delete(ctx context.Context, userId, walletId int) error {
	result, err := r.transactor.GetExecutor(ctx).ExecContext(ctx, deleteQuery, walletId, userId)
	if err != nil {
		return fmt.Errorf("[WalletPostgres.Delete] failed to delete wallet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[WalletPostgres.Delete] failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("[WalletPostgres.Delete] wallet not found")
	}
	return nil
}
