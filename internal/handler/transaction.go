package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goonsorrow/finance-tracker-api/internal/models"
)

type getAllTransactionsResponse struct {
	Wallet models.Wallet        `json:"wallet"`
	Data   []models.Transaction `json:"transactions"`
}

type getTransactionByIdResponse struct {
	Wallet models.Wallet      `json:"wallet"`
	Data   models.Transaction `json:"transaction"`
}

// @Summary Создать транзакцию
// @Description Пополнение (+) или списание (-) с кошелька
// @Security BearerAuth
// @Tags transactions
// @Accept json
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param input body models.CreateTransactionInput true "Сумма + Тип"
// @Success 200 {object} map[string]int "Transaction ID"
// @Router /api/wallets/{wallet_id}/transactions/ [post]
func (h *Handler) createTransaction(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	var input models.CreateTransactionInput
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

	id, err := h.services.Transaction.Create(ctx, userId, walletId, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while creating transaction")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Список транзакций кошелька
// @Description Получить все операции по кошельку
// @Security BearerAuth
// @Tags transactions
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Success 200 {object} handler.getAllTransactionsResponse
// @Router /api/wallets/{wallet_id}/transactions/ [get]
func (h *Handler) getAllTransactions(c *gin.Context) {
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

	items, err := h.services.Transaction.GetAll(ctx, userId, walletId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting transactions")
		return
	}

	c.JSON(http.StatusOK, getAllTransactionsResponse{
		Wallet: wallet,
		Data:   items,
	})
}

// @Summary Транзакция по ID
// @Security BearerAuth
// @Tags transactions
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param txId path int true "Transaction ID"
// @Success 200 {object} handler.getTransactionByIdResponse
// @Router /api/wallets/{wallet_id}/transactions/{txId} [get]
func (h *Handler) getTransactionByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id format")
		return
	}

	transactionId, err := strconv.Atoi(c.Param("txId"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid transaction id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	wallet, err := h.services.Wallet.GetById(ctx, userId, walletId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting wallet")
		return
	}

	transaction, err := h.services.Transaction.GetById(ctx, userId, walletId, transactionId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting transaction")
		return
	}

	c.JSON(http.StatusOK, getTransactionByIdResponse{
		Wallet: wallet,
		Data:   transaction,
	})
}

// @Summary Обновить транзакцию
// @Security BearerAuth
// @Tags transactions
// @Accept json
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param txId path int true "Transaction ID"
// @Param input body models.UpdateTransactionInput true "Changes"
// @Success 200 {object} handler.statusResponse
// @Router /api/wallets/{wallet_id}/transactions/{txId} [put]
func (h *Handler) updateTransactionByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id format")
		return
	}

	transactionId, err := strconv.Atoi(c.Param("txId"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid item id format")
		return
	}

	var input models.UpdateTransactionInput
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

	err = h.services.Transaction.Update(ctx, userId, walletId, transactionId, input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while updating todo item")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "transaction updated successfully",
	})
}

// @Summary Удалить транзакцию
// @Security BearerAuth
// @Tags transactions
// @Produce json
// @Param wallet_id path int true "Wallet ID"
// @Param txId path int true "Transaction ID"
// @Success 200 {object} handler.statusResponse
// @Router /api/wallets/{wallet_id}/transactions/{txId} [delete]
func (h *Handler) deleteTransactionByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	walletId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid wallet id format")
		return
	}

	transactionId, err := strconv.Atoi(c.Param("txId"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid transaction id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Transaction.Delete(ctx, userId, walletId, transactionId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while deleting transaction")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "transaction deleted successfully",
	})
}
