package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
)

func (s *MovementService) validateWalletAccess(ctx context.Context, userId, walletId int) error {
	_, err := s.walletRepo.GetById(ctx, userId, walletId)
	if err != nil {
		return err
	}
	return nil
}

type MovementService struct {
	walletRepo     repository.Wallet
	categoryRepo   repository.Category
	transactorRepo repository.Transactor
	movementRepo   repository.Movement
	logger         *slog.Logger
}

func NewMovementService(walletRepo repository.Wallet, categoryRepo repository.Category, transactorRepo repository.Transactor, movementRepo repository.Movement, logger *slog.Logger) *MovementService {
	return &MovementService{walletRepo: walletRepo, categoryRepo: categoryRepo, transactorRepo: transactorRepo, movementRepo: movementRepo, logger: logger}
}

func (s *MovementService) Create(ctx context.Context, userId, walletId int, input models.CreateMovementInput) (int, error) {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return 0, err
	}

	if err := input.Validate(); err != nil {
		s.logger.Error("validation error", slog.String("error", err.Error()))
		return 0, err
	}

	var movementId int

	amountInCents := int64(math.Round(input.Amount * 100))

	err := s.transactorRepo.WithinTransaction(ctx, func(txCtx context.Context) error {

		movement := models.Movement{
			WalletID:    walletId,
			UserId:      userId,
			Type:        input.Type,
			Amount:      amountInCents,
			CategoryID:  input.CategoryID,
			Description: input.Description,
			Date:        input.Date,
		}

		id, err := s.movementRepo.Create(txCtx, userId, walletId, movement)
		if err != nil {
			return fmt.Errorf("failed to create movement: %w", err)
		}
		movementId = id
		var diff int64

		if movement.Type == "expense" {
			diff = -amountInCents
		} else {
			diff = amountInCents
		}

		if diff != 0 {
			if err := s.walletRepo.AddToBalance(txCtx, walletId, diff); err != nil {
				return fmt.Errorf("failed to update wallet balance: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return movementId, nil
}

func (s *MovementService) CreateInitial(userId, walletId int, input models.CreateMovementInput) (int, error) {
	return 0, nil
}
func (s *MovementService) GetAll(ctx context.Context, userId, walletId int) ([]models.Movement, error) {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return []models.Movement{}, err
	}
	return s.movementRepo.GetAll(ctx, userId, walletId)
}

func (s *MovementService) GetById(ctx context.Context, userId, walletId, movementId int) (models.Movement, error) {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return models.Movement{}, err
	}
	return s.movementRepo.GetById(ctx, userId, walletId, movementId)
}

func (s *MovementService) Update(ctx context.Context, userId, walletId, movementId int, input models.UpdateMovementInput) error {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return err
	}

	if err := input.Validate(); err != nil {
		return err
	}

	err := s.transactorRepo.WithinTransaction(ctx, func(txCtx context.Context) error {
		oldMovement, err := s.movementRepo.GetById(txCtx, userId, walletId, movementId)
		if err != nil {
			return fmt.Errorf("failed to get old movement: %w", err)
		}

		oldDelta := oldMovement.Amount
		if oldMovement.Type == "expense" {
			oldDelta = -oldDelta
		}

		newAmount := oldMovement.Amount
		if input.Amount != nil {
			newAmount = int64(math.Round(*input.Amount * 100))
		}

		newType := oldMovement.Type
		if input.Type != nil {
			newType = *input.Type
		}

		newDelta := newAmount
		if newType == "expense" {
			newDelta = -newAmount
		}

		diff := newDelta - oldDelta

		if diff != 0 {
			if err := s.walletRepo.AddToBalance(txCtx, walletId, diff); err != nil {
				return fmt.Errorf("failed to update wallet balance: %w", err)
			}
		}

		updateInput := models.UpdateMovementData{
			Type:        input.Type,
			Amount:      nil,
			CategoryID:  input.CategoryID,
			Description: input.Description,
			Date:        input.Date,
		}
		if input.Amount != nil {
			amount := int64(math.Round(*input.Amount * 100))
			updateInput.Amount = &amount
		}

		if err := s.movementRepo.Update(txCtx, userId, walletId, movementId, updateInput); err != nil {
			return fmt.Errorf("failed to update movement: %w", err)
		}
		return nil
	})
	return err
}

func (s *MovementService) Delete(ctx context.Context, userId, walletId, movementId int) error {
	if err := s.validateWalletAccess(ctx, userId, walletId); err != nil {
		return err
	}

	err := s.transactorRepo.WithinTransaction(ctx, func(txCtx context.Context) error {
		oldMovement, err := s.movementRepo.GetById(txCtx, userId, walletId, movementId)
		if err != nil {
			return fmt.Errorf("failed to get old movement")
		}
		var delta int64
		if oldMovement.Type == "expense" {
			delta = oldMovement.Amount
		} else {
			delta = -oldMovement.Amount
		}

		if err := s.walletRepo.AddToBalance(txCtx, walletId, delta); err != nil {
			return fmt.Errorf("failed to update balance while deleting movement: %w", err)
		}

		if err := s.movementRepo.Delete(txCtx, userId, walletId, movementId); err != nil {
			return fmt.Errorf("failed to delete movement: %w", err)
		}
		return nil
	})
	return err
}
