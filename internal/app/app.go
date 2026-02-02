package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ultrabor/go-user-api/internal/config"
	"github.com/ultrabor/go-user-api/internal/delivery"
	"github.com/ultrabor/go-user-api/internal/domain"
	"github.com/ultrabor/go-user-api/internal/services"

	"github.com/ultrabor/go-user-api/internal/storage/memory"
	"github.com/ultrabor/go-user-api/internal/storage/postgres"

	_ "github.com/lib/pq"
)

func RunApp() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	usePostgres := true

	var store domain.UserStore

	if usePostgres {
		s, err := postgres.New(config.GetPostgresDSN(), logger)
		if err != nil {
			panic(err)
		}
		store = s
	} else {
		store = memory.New(logger)
	}

	userService := services.NewUserService(store)

	handler := delivery.NewUserHandler(userService)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(delivery.LoggingMiddleware(logger), gin.Recovery())

	handler.RegisterRoutes(router)

	srv := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		logger.Info("Server starting", "port", "8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Listen error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}
	logger.Info("Server exiting")
}
