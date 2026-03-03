package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
)

type WalletService struct {
	walletRepo      repository.Wallet
	movementRepo    repository.Movement
	logger          *slog.Logger
	transactor      repository.Transactor
	validCurrencies []string
}

func NewWalletService(walletRepo repository.Wallet, movementRepo repository.Movement, transactor repository.Transactor, logger *slog.Logger) *WalletService {

	return &WalletService{walletRepo: walletRepo, movementRepo: movementRepo, logger: logger, transactor: transactor, validCurrencies: []string{"USD", "EUR", "RUB", "GBP", "JPY"}}
}

func (s *WalletService) Create(ctx context.Context, userId int, input models.CreateWalletInput) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, err
	}
	if err := s.ValidateCurrency(ctx, input.Currency); err != nil {
		return 0, err
	}
	balanceInCents := int64(math.Round(input.InitialBalance * 100))

	var walletId int

	err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error

		wallet := models.Wallet{
			UserID:   userId,
			Name:     input.Name,
			Currency: input.Currency,
			Balance:  balanceInCents, // Пишем копейки
		}

		walletId, err = s.walletRepo.Create(txCtx, userId, wallet)
		if err != nil {
			return err
		}

		if balanceInCents != 0 {
			initialMovement := models.Movement{
				WalletID:    walletId,
				UserId:      userId,
				Type:        "initial",
				Amount:      balanceInCents,
				CategoryID:  nil,
				Description: "Initial balance set at wallet creation",
				Date:        time.Now(),
			}
			_, err = s.movementRepo.Create(txCtx, userId, walletId, initialMovement)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return walletId, nil
}

func (s *WalletService) GetAll(ctx context.Context, userId int) ([]models.Wallet, error) {
	return s.walletRepo.GetAll(ctx, userId)
}

func (s *WalletService) GetById(ctx context.Context, userId, walletId int) (models.Wallet, error) {
	return s.walletRepo.GetById(ctx, userId, walletId)
}

func (s *WalletService) Update(ctx context.Context, userId, walletId int, input models.UpdateWalletInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.walletRepo.Update(ctx, userId, walletId, input)
}

func (s *WalletService) Delete(ctx context.Context, userId, walletId int) error {
	return s.walletRepo.Delete(ctx, userId, walletId)
}

func (s *WalletService) ValidateCurrency(ctx context.Context, currency string) error {
	for _, valid := range s.validCurrencies {
		if currency == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid currency: %s valid: %v", currency, s.validCurrencies)
}
