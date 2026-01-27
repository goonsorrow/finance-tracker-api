package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
)

type WalletService struct {
	repo            repository.Wallet
	trRepo          repository.Transaction
	logger          *slog.Logger
	validCurrencies []string
}

func NewWalletService(walletRepo repository.Wallet, logger *slog.Logger) *WalletService {

	return &WalletService{repo: walletRepo, logger: logger, validCurrencies: []string{"USD", "EUR", "RUB", "GBP", "JPY"}}
}

func (s *WalletService) Create(ctx context.Context, userId int, input models.CreateWalletInput) (int, error) {
	if err := input.Validate(); err != nil {
		return 0, err
	}
	if err := s.ValidateCurrency(ctx, input.Currency); err != nil {
		return 0, err
	}
	balanceInCents := int64(input.InitialBalance * 100)

	wallet := models.Wallet{
		UserID:   userId,
		Name:     input.Name,
		Currency: input.Currency,
		Balance:  balanceInCents, // Пишем копейки
	}

	if balanceInCents > 0 {
		return s.repo.CreateWithInitial(ctx, userId, wallet)
	}

	return s.repo.Create(ctx, userId, wallet)
}

func (s *WalletService) GetAll(ctx context.Context, userId int) ([]models.Wallet, error) {
	return s.repo.GetAll(ctx, userId)
}

func (s *WalletService) GetById(ctx context.Context, userId, walletId int) (models.Wallet, error) {
	return s.repo.GetById(ctx, userId, walletId)
}

func (s *WalletService) Update(ctx context.Context, userId, walletId int, input models.UpdateWalletInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	updateInput := models.UpdateWalletData{
		Name: input.Name,
	}

	return s.repo.Update(ctx, userId, walletId, updateInput)
}

func (s *WalletService) Delete(ctx context.Context, userId, walletId int) error {
	return s.repo.Delete(ctx, userId, walletId)
}

func (s *WalletService) ValidateCurrency(ctx context.Context, currency string) error {
	for _, valid := range s.validCurrencies {
		if currency == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid currency: %s valid: %v", currency, s.validCurrencies)
}
