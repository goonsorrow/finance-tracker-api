package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/goonsorrow/finance-tracker/internal/service"
)

func NewHandler(services *service.Service, logger *slog.Logger) *Handler {
	return &Handler{services: services, logger: logger}
}

func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, err error, message string) {
	if statusCode >= 500 {
		h.logger.Error(message, slog.String("error", err.Error()))
	} else {
		h.logger.Warn(message, slog.String("error", err.Error()))
	}

	c.AbortWithStatusJSON(statusCode, map[string]interface{}{
		"error": message,
	})
}
