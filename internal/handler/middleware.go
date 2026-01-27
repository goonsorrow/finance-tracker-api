package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	userCtx             = "userId"
	autorizathionHeader = "Authorization"
)

// @SecurityDefinitions.apikey Bearer
// @In header
// @Name Authorization
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(autorizathionHeader)
	if header == "" {
		h.newErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("empty header"), "empty header")
		return
	}

	headerSlice := strings.Split(header, " ")
	if len(headerSlice) != 2 {
		h.newErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("invalid format"), "invalid format")
		return
	}

	if headerSlice[0] != "Bearer" {
		h.newErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("not Bearer"), "not Bearer")
		return
	}

	token := headerSlice[1]

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userId, err := h.services.Authorization.
		ParseAccessToken(ctx, token)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("ParseAccessToken Failed"), err.Error())
		return
	}
	h.logger.Info("auth middleware passed",
		slog.Int("user_id", userId))

	c.Set(userCtx, userId)
	c.Next()
}

func (h *Handler) getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		h.newErrorResponse(c, http.StatusInternalServerError, nil, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		h.newErrorResponse(c, http.StatusInternalServerError, nil, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil

}

func (h *Handler) LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		userId, exists := ctx.Get(userCtx)
		userAttr := slog.String("user_id", "guest")
		if exists {
			if id, ok := userId.(int); ok {
				userAttr = slog.Int("user_id", id)
			}
		}

		attrs := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("clientIP", clientIP),
			userAttr,
		}

		if status >= 500 {
			h.logger.Error("Request failed", attrs...)
		} else if status >= 400 {
			h.logger.Warn("Bad request", attrs...)
		} else {
			h.logger.Info("Request success", attrs...)
		}
	}
}
