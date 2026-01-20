package repository

import (
	"context"

	"github.com/goonsorrow/finance-tracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserById(ctx context.Context, id int) (models.User, error)
	CreateRefreshSession(ctx context.Context, session models.RefreshSession) error
	GetRefreshSession(ctx context.Context, token string) (models.RefreshSession, error)
	DeleteRefreshSession(ctx context.Context, token string) error
}
type Wallet interface {
	Create(ctx context.Context, userId int, wallet models.Wallet) (int, error)
	CreateWithInitial(ctx context.Context, userId int, wallet models.Wallet) (int, error)
	GetAll(ctx context.Context, userId int) ([]models.Wallet, error)
	GetById(ctx context.Context, userId, walletId int) (models.Wallet, error)
	Update(ctx context.Context, userId, walletId int, input models.UpdateWalletData) error
	Delete(ctx context.Context, userId, walletId int) error
}
type Transaction interface {
	Create(ctx context.Context, userId, walletId int, transaction models.Transaction) (int, error)
	GetAll(ctx context.Context, userId, walletId int) ([]models.Transaction, error)
	GetById(ctx context.Context, userId, walletId, transactionId int) (models.Transaction, error)
	Delete(ctx context.Context, userId, walletId, transactionId int) error
	Update(ctx context.Context, userId, walletId, transactionId int, input models.UpdateTransactionData) error
}

type Repository struct {
	Authorization
	Wallet
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Wallet:        NewWalletPostgres(db),
		Transaction:   NewTransactionPostgres(db),
	}
}
