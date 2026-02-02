package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/goonsorrow/finance-tracker-api/configs"
	"github.com/goonsorrow/finance-tracker-api/internal/cache"
	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo       repository.Authorization
	cache      cache.Authorization
	logger     *slog.Logger
	jwtConfig  configs.JWTConfig
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo repository.Authorization, cache cache.Authorization, logger *slog.Logger, jwtConfig configs.JWTConfig) *AuthService {
	accessTTL, err := time.ParseDuration(jwtConfig.AccessTTL)
	if err != nil {
		logger.Warn("invalid access_ttl config, using default 15m", "error", err)
		accessTTL = 15 * time.Minute
	}

	refreshTTL, err := time.ParseDuration(jwtConfig.RefreshTTL)
	if err != nil {
		logger.Warn("invalid refresh_ttl config, using default 24h", "error", err)
		refreshTTL = 24 * time.Hour
	}
	return &AuthService{repo: repo, cache: cache, logger: logger, jwtConfig: jwtConfig, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (s *AuthService) CreateUser(ctx context.Context, input models.RegisterInput) (int, error) {
	hashedPassword, err := generatePasswordHash(input.Password)
	if err != nil {
		return 0, err
	}
	user := models.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to find user by email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.logger.Warn("invalid password attempt", slog.String("email", email))
		return "", "", errors.New("wrong password")
	}

	return s.createSession(ctx, user.ID, user.Email)
}

func (s *AuthService) createSession(ctx context.Context, userId int, email string) (string, string, error) {
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, &AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email, // Храним email здесь, чтобы потом достать при рефреше
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		UserId: userId,
		Email:  email,
	})
	accessToken, err := accessTokenObj.SignedString([]byte(s.jwtConfig.SigningKey))
	if err != nil {
		return "", "", fmt.Errorf("sign access token error: %w", err)
	}

	refreshExpiresAt := time.Now().UTC().Add(s.refreshTTL)
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		UserId: userId,
	})
	refreshToken, err := refreshTokenObj.SignedString([]byte(s.jwtConfig.SigningKey))
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token error: %w", err)
	}

	session := models.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: refreshExpiresAt,
	}

	key := fmt.Sprintf("refresh:userId:%d:%s", userId, refreshToken)
	if err := s.cache.CacheRefreshSession(ctx, key, s.refreshTTL); err != nil {
		return "", "", fmt.Errorf("cache save error: %w", err)
	}

	if err := s.repo.CreateRefreshSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("db save error: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, oldRefreshToken string) (string, string, error) {
	claims, err := s.ValidateRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token presented: %w", err)
	}

	session, err := s.repo.GetRefreshSession(ctx, oldRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("refresh token not found in db: %w", err)
	}

	if session.ExpiresAt.Before(time.Now().UTC()) {
		_ = s.repo.DeleteRefreshSession(ctx, oldRefreshToken)
		return "", "", errors.New("refresh token is expired")
	}

	if err := s.repo.DeleteRefreshSession(ctx, oldRefreshToken); err != nil {
		s.logger.Error("failed to delete old session", "error", err)
	}
	user, err := s.repo.GetUserById(ctx, claims.UserId)
	if err != nil {
		s.logger.Error("failed to find user", "error", err)
		return "", "", err
	}

	return s.createSession(ctx, user.ID, user.Email)

}

func (s *AuthService) ParseAccessToken(ctx context.Context, accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.jwtConfig.SigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)

	if !ok {
		return 0, errors.New("accessToken claims are not of type *AccessTokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) (string, error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", errors.New("error while hashing password")
	}

	return string(hash), nil
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.jwtConfig.SigningKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, err
	}

	key := fmt.Sprintf("refresh:userId:%d:%s", claims.UserId, refreshToken)
	result, err := s.cache.CheckRefreshToken(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("redis error:%w", err)
	}
	if result == 0 {
		return nil, fmt.Errorf("refresh token has expired")
	}

	return claims, nil

}
