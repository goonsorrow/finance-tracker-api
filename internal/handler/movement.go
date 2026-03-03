package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goonsorrow/finance-tracker-api/internal/models"
)

type getAllMovementsResponse struct {
	Wallet models.Wallet     `json:"wallet"`
	Data   []models.Movement `json:"movements"`
}

type getMovementByIdResponse struct {
	Wallet models.Wallet   `json:"wallet"`
	Data   models.Movement `json:"movement"`
}

// @Summary Создать транзакцию
// @Description Пополнение (+) или списание (-) с кошелька
// @Security Bearer
// @Tags movements
// @Accept json
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param input body models.CreateMovementInput true "Сумма + Тип"
// @Success 200 {object} map[string]int "Movement ID"
// @Router /api/wallets/{wallet_id}/movements/ [post]
func (h *Handler) createMovement(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	var input models.CreateMovementInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input data")
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := h.services.Movement.Create(ctx, userId, walletId, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while creating movement")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Список транзакций кошелька
// @Description Получить все операции по кошельку
// @Security Bearer
// @Tags movements
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Success 200 {object} handler.getAllMovementsResponse
// @Router /api/wallets/{wallet_id}/movements/ [get]
func (h *Handler) getAllMovements(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	wallet, err := h.services.Wallet.GetById(ctx, userId, walletId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting wallet")
		return
	}

	items, err := h.services.Movement.GetAll(ctx, userId, walletId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting movements")
		return
	}

	c.JSON(http.StatusOK, getAllMovementsResponse{
		Wallet: wallet,
		Data:   items,
	})
}

// @Summary Транзакция по ID
// @Security Bearer
// @Tags movements
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param trId path int true "Movement ID"
// @Success 200 {object} handler.getMovementByIdResponse
// @Router /api/wallets/{wallet_id}/movements/{trId} [get]
func (h *Handler) getMovementByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id format")
		return
	}

	movementId, err := strconv.Atoi(c.Param("trId"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid movement id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	wallet, err := h.services.Wallet.GetById(ctx, userId, walletId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting wallet")
		return
	}

	movement, err := h.services.Movement.GetById(ctx, userId, walletId, movementId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting movement")
		return
	}

	c.JSON(http.StatusOK, getMovementByIdResponse{
		Wallet: wallet,
		Data:   movement,
	})
}

// @Summary Обновить транзакцию
// @Security Bearer
// @Tags movements
// @Accept json
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param trId path int true "Movement ID"
// @Param input body models.UpdateMovementInput true "Changes"
// @Success 200 {object} handler.statusResponse
// @Router /api/wallets/{wallet_id}/movements/{trId} [put]
func (h *Handler) updateMovementByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id format")
		return
	}

	movementId, err := strconv.Atoi(c.Param("trId"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid item id format")
		return
	}

	var input models.UpdateMovementInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input data")
		return
	}

	if err := input.Validate(); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Movement.Update(ctx, userId, walletId, movementId, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while updating todo item")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "movement updated successfully",
	})
}

// @Summary Удалить транзакцию
// @Security Bearer
// @Tags movements
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param trId path int true "Movement ID"
// @Success 200 {object} handler.statusResponse
// @Router /api/wallets/{wallet_id}/movements/{trId} [delete]
func (h *Handler) deleteMovementByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id format")
		return
	}

	movementId, err := strconv.Atoi(c.Param("trId"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid movement id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Movement.Delete(ctx, userId, walletId, movementId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while deleting movement")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "movement deleted successfully",
	})
}
