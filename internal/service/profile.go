package service

import (
	"context"
	"log/slog"
)

type ProfileService struct {
	logger *slog.Logger
}

func NewProfileService(logger *slog.Logger) *ProfileService {
	return &ProfileService{logger: logger}
}

func (ps ProfileService) GetMe(ctx context.Context, userId int) {
	// we need to get username, profile and his balance in all currencies divided by them or combined into one by exchange rate.

}
