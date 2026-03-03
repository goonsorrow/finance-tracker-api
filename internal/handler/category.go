package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goonsorrow/finance-tracker-api/internal/models"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
)

type getAllCategoriesResponse struct {
	Data []models.Category `json:"categories"`
}

// @Summary Создать категорию
// @Description Добавить новую категорию
// @Security Bearer
// @Tags categories
// @Accept json
// @Produce json
// @Param input body models.CreateCategoryInput true "Name + Type + Icon"
// @Success 201 {object} map[string]int
// @Router /api/categories/ [post]
func (h *Handler) createCategory(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	var input models.CreateCategoryInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid input data")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := h.services.Category.Create(ctx, userId, input)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			h.newErrorResponse(c, http.StatusNotFound, err, "category not found")
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while creating category")
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// @Summary Список категорий
// @Description Получить все категории пользователя
// @Security Bearer
// @Tags categories
// @Produce json
// @Success 200 {object} handler.getAllCategoriesResponse
// @Failure 401 {object} map[string]string
// @Router /api/categories [get]
func (h *Handler) getAllCategories(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	categories, err := h.services.Category.GetAll(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			h.newErrorResponse(c, http.StatusNotFound, err, "category not found")
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, err, "failed to get categories")
		return
	}

	c.JSON(http.StatusOK, getAllCategoriesResponse{
		Data: categories,
	})
}

// @Summary Получить категорию по ID
// @Description Детали конкретной категории
// @Security Bearer
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Not found"
// @Router /api/categories/{id} [get]
func (h *Handler) getCategoryByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid category id")
		return
	}
	category, err := h.services.Category.GetById(ctx, userId, categoryId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			h.newErrorResponse(c, http.StatusNotFound, err, "category not found")
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while getting user category by id")
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Получить категорию по ID
// @Description Детали конкретной категории
// @Security Bearer
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Not found"
// @Router /api/categories/{id} [get]
func (h *Handler) updateCategoryByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid category id")
		return
	}

	var input models.UpdateCategoryInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "error while reading input")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Category.Update(ctx, userId, categoryId, input)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			h.newErrorResponse(c, http.StatusNotFound, err, "category not found")
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while updating user category by id")
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// @Summary Удалить категорию
// @Description Удалить категорию
// @Security Bearer
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} handler.statusResponse
// @Failure 400 {object} map[string]string "Invalid ID"
// @Router /api/categories/{id} [delete]
func (h *Handler) deleteCategoryByID(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err, "invalid category id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.services.Category.Delete(ctx, userId, categoryId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			h.newErrorResponse(c, http.StatusNotFound, err, "category not found")
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, err, "error while deleting user category by id")
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
