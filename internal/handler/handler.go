// @title Finance Tracker API
// @version 1.0
// @description Go REST API для финансового трекера (JWT + Docker)
// @host localhost:8080
// @BasePath /

package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	_ "github.com/goonsorrow/finance-tracker-api/docs"
	"github.com/goonsorrow/finance-tracker-api/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	logger   *slog.Logger
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.LoggingMiddleware())
	router.Use(gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.POST("/login", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout-all", h.logoutAll)
		auth.POST("/logout", h.logout)
	}
	api := router.Group("/api")
	api.Use(h.userIdentity)

	wallets := api.Group("/wallets")
	{
		wallets.GET("/", h.getAllWallets)
		wallets.GET("/:id", h.getWalletByID)
		wallets.POST("/", h.createWallet)
		wallets.PUT("/:id", h.updateWalletByID)
		wallets.DELETE("/:id", h.deleteWalletByID)

		transactions := wallets.Group("/:id/transactions")
		{
			transactions.GET("/", h.getAllTransactions)
			transactions.GET("/:trId", h.getTransactionByID)
			transactions.POST("/", h.createTransaction)
			transactions.PUT("/:trId", h.updateTransactionByID)
			transactions.DELETE("/:trId", h.deleteTransactionByID)
		}

	}
	categories := api.Group("/categories")
	{
		categories.GET("/", h.getAllCategories)
		categories.GET("/:id", h.getCategoryByID)
		categories.POST("/", h.createCategory)
		categories.PUT("/:id", h.updateCategoryByID)
		categories.DELETE("/:id", h.deleteCategoryByID)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
