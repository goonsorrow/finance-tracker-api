package service

import (
	"context"
	"log/slog"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
)

type CategoryService struct {
	repo   repository.Category
	logger *slog.Logger
}

func NewCategoryService(categoryRepo repository.Category, logger *slog.Logger) *CategoryService {

	return &CategoryService{repo: categoryRepo, logger: logger}
}

func (s *CategoryService) Create(ctx context.Context, userId int, input models.CreateCategoryInput) (int, error) {
	category := models.Category{
		Name:   input.Name,
		Type:   input.Type,
		Icon:   input.Icon,
		UserID: &userId,
	}

	return s.repo.Create(ctx, userId, category)
}

func (s *CategoryService) GetAll(ctx context.Context, userId int) ([]models.Category, error) {
	return s.repo.GetAll(ctx, userId)
}

func (s *CategoryService) GetById(ctx context.Context, userId, categoryId int) (models.Category, error) {
	return s.repo.GetById(ctx, userId, categoryId)
}

func (s *CategoryService) Update(ctx context.Context, userId, categoryId int, input models.UpdateCategoryInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, userId, categoryId, input)
}

func (s *CategoryService) Delete(ctx context.Context, userId, categoryId int) error {
	return s.repo.Delete(ctx, userId, categoryId)
}
