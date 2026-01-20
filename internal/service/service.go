package service

import (
	"context"
	"log/slog"

	"github.com/goonsorrow/finance-tracker/configs"
	"github.com/goonsorrow/finance-tracker/internal/models"
	"github.com/goonsorrow/finance-tracker/internal/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.RegisterInput) (int, error)
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	ParseAccessToken(ctx context.Context, token string) (int, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenClaims, error)
	RefreshTokens(ctx context.Context, oldRefreshToken string) (string, string, error)
	createSession(ctx context.Context, userId int, email string) (string, string, error)
}
type Wallet interface {
	Create(ctx context.Context, userId int, wallet models.CreateWalletInput) (int, error)
	GetAll(ctx context.Context, userId int) ([]models.Wallet, error)
	GetById(ctx context.Context, userId, walletId int) (models.Wallet, error)
	Delete(ctx context.Context, userId, walletId int) error
	Update(ctx context.Context, userId, walletId int, input models.UpdateWalletInput) error
}
type Transaction interface {
	Create(ctx context.Context, userId int, walletId int, transaction models.CreateTransactionInput) (int, error)
	GetAll(ctx context.Context, userId, walletId int) ([]models.Transaction, error)
	GetById(ctx context.Context, userId, walletId, transactionId int) (models.Transaction, error)
	Delete(ctx context.Context, userId, walletId, transactionId int) error
	Update(ctx context.Context, userId, walletId, transactionId int, input models.UpdateTransactionInput) error
}

type Service struct {
	Authorization
	Wallet
	Transaction
	logger *slog.Logger
}

func NewService(repos *repository.Repository, logger *slog.Logger, cfg configs.Config) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, logger, cfg.JWT),
		Wallet:        NewWalletService(repos.Wallet, logger),
		Transaction:   NewTransactionService(repos.Wallet, repos.Transaction, logger),
		logger:        logger,
	}
}
