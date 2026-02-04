package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goonsorrow/finance-tracker-api/internal/models"
)

// @Summary Регистрация пользователя
// @Description Создать аккаунт
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RegisterInput true "Email + Password"
// @Success 201 {object} map[string]int "User ID"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /auth/register [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := h.services.Authorization.CreateUser(ctx, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "failed to register")
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// @Summary Вход
// @Description Аутентификация и получение JWT токенов
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body models.SignInInput true "Credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /auth/login [post]
func (h *Handler) signIn(c *gin.Context) {
	var input models.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accessToken, refreshToken, err := h.services.Authorization.SignIn(ctx, input.Email, input.Password)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary Обновить токены
// @Description Refresh access по refresh токену
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RefreshInput true "Refresh token"
// @Success 200 {object} map[string]string "New tokens"
// @Failure 401 {object} map[string]string "Invalid token"
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input models.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accessToken, refreshToken, err := h.services.Authorization.RefreshTokens(ctx, input.RefreshToken)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err, "invalid token")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary Выйти со конкретного устройства
// @Description Отозвать конкретную refresh сессию пользователя (logout here)
// @Tags auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.LogoutInput true "Refresh token"
// @Success 200 {object} map[string]string "Session deleted"
// @Failure 401 {object} map[string]string "Not Authorized"
// @Failure 500 {object} map[string]string "Failed deleting session"
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	var input models.LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("logout failed:%w", err), "logout fauled - invalid request")
		return
	}
	claims, err := h.services.ValidateRefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("logout failed:%w", err), "failed validating refresh token")
		return
	}

	if err := h.services.Authorization.LogoutCurrentUserSession(c.Request.Context(), claims.UserId, claims.ID); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("logout failed:%w", err), "logout failed")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "deleted current session"})
}

// @Summary Выйти со всех устройств
// @Description Отозвать все активные refresh сессии пользователя (logout everywhere)
// @Tags auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.LogoutInput true "Refresh token"
// @Success 200 {object} map[string]string "All sessions deleted"
// @Failure 401 {object} map[string]string "Not Authorized"
// @Failure 500 {object} map[string]string "Failed deleting sessions"
// @Router /auth/logout-all [post]
func (h *Handler) logoutAll(c *gin.Context) {
	var input models.LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, fmt.Errorf("logout failed:%w", err), "logout fauled - invalid request")
		return
	}
	claims, err := h.services.ValidateRefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("logout failed:%w", err), "failed validating refresh token")
		return
	}

	if err := h.services.Authorization.LogoutAllUserSessions(c.Request.Context(), claims.UserId); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("logout failed:%w", err), "logout failed")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "all sessions logged out"})
}
