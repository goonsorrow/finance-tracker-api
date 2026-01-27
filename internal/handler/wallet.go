package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goonsorrow/finance-tracker-api/internal/models"
)

type getAllWalletsResponse struct {
	Data []models.Wallet `json:"wallets"`
}

type statusResponse struct {
	Status string `json:"status"`
}

// @Summary Создать кошелёк
// @Description Добавить новый кошелёк
// @Security Bearer
// @Tags wallets
// @Accept json
// @Produce json
// @Param input body models.CreateWalletInput true "Name + Currency"
// @Success 201 {object} map[string]int
// @Router /api/wallets/ [post]
func (h *Handler) createWallet(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	var input models.CreateWalletInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input data")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := h.services.Wallet.Create(ctx, userId, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while creating wallet")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Список кошельков
// @Description Получить все кошельки пользователя
// @Security Bearer
// @Tags wallets
// @Produce json
// @Success 200 {object} handler.getAllWalletsResponse
// @Failure 401 {object} map[string]string
// @Router /api/wallets/ [get]
func (h *Handler) getAllWallets(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	wallets, err := h.services.Wallet.GetAll(ctx, userId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "failed to get wallets")
		return
	}

	c.JSON(http.StatusOK, getAllWalletsResponse{
		Data: wallets,
	})
}

// @Summary Получить кошелёк по ID
// @Description Детали конкретного кошелька
// @Security Bearer
// @Tags wallets
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Not found"
// @Router /api/wallets/{id} [get]
func (h *Handler) getWalletByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id")
		return
	}
	wallet, err := h.services.Wallet.GetById(ctx, userId, walletId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting user wallet by id")
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// @Summary Получить кошелёк по ID
// @Description Детали конкретного кошелька
// @Security Bearer
// @Tags wallets
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Not found"
// @Router /api/wallets/{id} [get]
func (h *Handler) updateWalletByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "user not found")
		return
	}

	var input models.UpdateWalletInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "error while reading input")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Wallet.Update(ctx, userId, id, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while updating user wallet by id")
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// @Summary Удалить кошелёк
// @Description Удалить кошелёк и все транзакции
// @Security Bearer
// @Tags wallets
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {object} handler.statusResponse
// @Failure 400 {object} map[string]string "Invalid ID"
// @Router /api/wallets/{id} [delete]
func (h *Handler) deleteWalletByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "user not found")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Wallet.Delete(ctx, userId, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while deleting user wallet by id")
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
