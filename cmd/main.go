package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/goonsorrow/finance-tracker-api/configs"
	"github.com/goonsorrow/finance-tracker-api/internal/app"
	"github.com/goonsorrow/finance-tracker-api/internal/handler"
	"github.com/goonsorrow/finance-tracker-api/internal/logger"
	"github.com/goonsorrow/finance-tracker-api/internal/repository"
	"github.com/goonsorrow/finance-tracker-api/internal/service"
	"github.com/spf13/viper"
)

// @title Finance Tracker API
// @version 1.0
// @description Personal finance tracker API. JWT + Postgres + Docker.
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slogger := logger.InitSlogger()

	if err := InitConfig(); err != nil {
		slogger.Error("error occured while initialising config:", "error", err)
		os.Exit(1)
	}

	var cfg configs.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		slogger.Error("error occured while reading config:", "error", err)
		os.Exit(1)
	}

	if cfg.JWT.SigningKey == "" {
		slogger.Error("JWT signing key is not set (jwt.signing_key in config)")
		os.Exit(1)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		slogger.Error("error occured while connecting to db:", "err", err)
		os.Exit(1)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, 6379),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Проверяем подключение
	_, err = client.Ping(ctx).Result()
	if err != nil {
		slogger.Error("error occured while connecting to redis:", "err", err)
		return
	}

	repo := repository.NewRepository(db)
	// cache := cache.NewCache(client)
	service := service.NewService(repo, slogger, cfg)
	handler := handler.NewHandler(service, slogger)
	srv := new(app.Server)

	go func() {
		if err := srv.Run(cfg.Server.Port, handler.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogger.Error("Error occured while running http server", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slogger.Info("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slogger.Error("Error while gracefully shutting down server", "error", err)
	}
	if err := db.Close(); err != nil {
		slogger.Error("Error while closing db", "error", err)
		os.Exit(1)
	}

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.BindEnv("jwt.access_ttl", "JWT_ACCESS_TTL")
	_ = viper.BindEnv("jwt.refresh_ttl", "JWT_REFRESH_TTL")
	_ = viper.BindEnv("jwt.signing_key", "JWT_SIGNING_KEY")

	// Server
	_ = viper.BindEnv("server.port", "SERVER_PORT")

	// DB
	_ = viper.BindEnv("db.host", "DB_HOST")
	_ = viper.BindEnv("db.port", "DB_PORT")
	_ = viper.BindEnv("db.username", "DB_USERNAME")
	_ = viper.BindEnv("db.password", "DB_PASSWORD")
	_ = viper.BindEnv("db.dbname", "DB_DBNAME")
	_ = viper.BindEnv("db.sslmode", "DB_SSLMODE")
	// Redis
	_ = viper.BindEnv("redis.host", "REDIS_HOST")
	_ = viper.BindEnv("redis.port", "REDIS_PORT")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("db.sslmode", "disable")
	return nil
}
