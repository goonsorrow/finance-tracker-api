package service

import (
	"context"
	"log/slog"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
)

func (s *TransactionService) validateWalletAccess(ctx context.Context, userId, walletId int) error {
	_, err := s.walletRepo.GetById(ctx, userId, walletId)
	if err != nil {
		return err
	}
	return nil
}

type TransactionService struct {
	walletRepo      repository.Wallet
	transactionRepo repository.Transaction
	logger          *slog.Logger
}

func NewTransactionService(walletRepo repository.Wallet, transactionRepo repository.Transaction, logger *slog.Logger) *TransactionService {
	return &TransactionService{walletRepo: walletRepo, transactionRepo: transactionRepo, logger: logger}
}

func (s *TransactionService) Create(ctx context.Context, userId, walletId int, input models.CreateTransactionInput) (int, error) {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return 0, err
	}

	if err := input.Validate(); err != nil {
		s.logger.Error("validation error", slog.String("error", err.Error()))
		return 0, err
	}
	amountInCents := int64(input.Amount * 100)

	transaction := models.Transaction{
		WalletID:    walletId,
		UserId:      userId,
		Type:        input.Type,
		Amount:      amountInCents,
		Category:    input.Category,
		Description: input.Description,
		Date:        input.Date,
	}

	return s.transactionRepo.Create(ctx, userId, walletId, transaction)
}

func (s *TransactionService) CreateInitial(userId, walletId int, input models.CreateTransactionInput) (int, error) {
	return 0, nil
}
func (s *TransactionService) GetAll(ctx context.Context, userId, walletId int) ([]models.Transaction, error) {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return []models.Transaction{}, err
	}
	return s.transactionRepo.GetAll(ctx, userId, walletId)
}

func (s *TransactionService) GetById(ctx context.Context, userId, walletId, transactionId int) (models.Transaction, error) {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return models.Transaction{}, err
	}
	return s.transactionRepo.GetById(ctx, userId, walletId, transactionId)
}

func (s *TransactionService) Update(ctx context.Context, userId, walletId, transactionId int, input models.UpdateTransactionInput) error {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return err
	}

	if err := input.Validate(); err != nil {
		return err
	}

	var amountInCents *int64
	if input.Amount != nil {
		a := int64(*input.Amount * 100)
		amountInCents = &a
	}

	updateInput := models.UpdateTransactionData{
		Type:        input.Type,
		Amount:      amountInCents,
		Category:    input.Category,
		Description: input.Description,
		Date:        input.Date,
	}

	return s.transactionRepo.Update(ctx, userId, walletId, transactionId, updateInput)
}

func (s *TransactionService) Delete(ctx context.Context, userId, walletId, transactionId int) error {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return err
	}
	return s.transactionRepo.Delete(ctx, userId, walletId, transactionId)
}
