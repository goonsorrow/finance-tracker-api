package repository

import (
	"context"
	"fmt"

	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

const (
	createUserQuery    = `INSERT INTO users (email,password_hash,created_at,updated_at) VALUES ($1,$2,NOW(),NOW()) RETURNING id`
	getUserByEmail     = `SELECT id, email, password_hash,created_at,updated_at FROM users WHERE email=$1`
	getUserByIdQuery   = `SELECT id, email, password_hash,created_at,updated_at FROM users WHERE id=$1`
	createSessionQuery = `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	getSessionQuery    = `SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = $1`
	deleteSessionQuery = `DELETE FROM refresh_tokens WHERE token=$1`
)

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user models.User) (int, error) {
	var id int
	row := r.db.QueryRowContext(ctx, createUserQuery,
		user.Email,        //$1
		user.PasswordHash) //$2
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, getUserByEmail, email)
	if err != nil {
		return user, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *AuthPostgres) GetUserById(ctx context.Context, id int) (models.User, error) {
	var user models.User

	err := r.db.GetContext(ctx, &user, getUserByIdQuery, id)
	if err != nil {
		return user, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *AuthPostgres) CreateRefreshSession(ctx context.Context, s models.RefreshSession) error {
	_, err := r.db.ExecContext(ctx, createSessionQuery,
		s.UserID,    //$1
		s.Token,     //$2
		s.ExpiresAt, //$3
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh session: %w", err)
	}
	return nil
}

func (r *AuthPostgres) GetRefreshSession(ctx context.Context, token string) (models.RefreshSession, error) {
	var session models.RefreshSession
	err := r.db.GetContext(ctx, &session, getSessionQuery, token)
	if err != nil {
		return models.RefreshSession{}, fmt.Errorf("failed to get refresh session")
	}
	return session, nil
}

func (r *AuthPostgres) DeleteRefreshSession(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, deleteSessionQuery, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh session:%w", err)
	}
	return nil
}
