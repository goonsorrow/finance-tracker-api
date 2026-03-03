package repository

import (
	"context"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	GetExecutor(ctx context.Context) sqlx.ExtContext
}

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
	GetAll(ctx context.Context, userId int) ([]models.Wallet, error)
	GetById(ctx context.Context, userId, walletId int) (models.Wallet, error)
	Update(ctx context.Context, userId, walletId int, input models.UpdateWalletInput) error
	AddToBalance(ctx context.Context, walletId int, deltaCents int64) error
	Delete(ctx context.Context, userId, walletId int) error
}
type Movement interface {
	Create(ctx context.Context, userId, walletId int, movement models.Movement) (int, error)
	GetAll(ctx context.Context, userId, walletId int) ([]models.Movement, error)
	GetById(ctx context.Context, userId, walletId, movementId int) (models.Movement, error)
	Update(ctx context.Context, userId, walletId, movementId int, input models.UpdateMovementData) error
	Delete(ctx context.Context, userId, walletId, movementId int) error
}

type Category interface {
	Create(ctx context.Context, userId int, category models.Category) (int, error)
	GetAll(ctx context.Context, userId int) ([]models.Category, error)
	GetById(ctx context.Context, userId, categoryId int) (models.Category, error)
	Update(ctx context.Context, userId, categoryId int, input models.UpdateCategoryInput) error
	Delete(ctx context.Context, userId, categoryId int) error
}

type Repository struct {
	Transactor
	Authorization
	Wallet
	Movement
	Category
}

func NewRepository(db *sqlx.DB) *Repository {
	transactor := NewTransactorPostgres(db)
	return &Repository{
		Transactor:    transactor,
		Authorization: NewAuthPostgres(db),
		Wallet:        NewWalletPostgres(db, transactor),
		Movement:      NewMovementPostgres(db, transactor),
		Category:      NewCategoryPostgres(db, transactor),
	}
}
