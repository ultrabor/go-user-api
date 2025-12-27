package server

import (
	"log/slog"
	"net/http"

	"github.com/ultrabor/go-user-api/internal/handlers"
	"github.com/ultrabor/go-user-api/internal/middleware"
	"github.com/ultrabor/go-user-api/internal/services"
	"github.com/ultrabor/go-user-api/internal/storage"
)

func Server(logger *slog.Logger, store storage.UserStore) http.Handler {

	mux := http.NewServeMux()

	userService := services.NewUserService(store)

	mux.HandleFunc("/create", handlers.CreateUserHandler(userService))
	mux.HandleFunc("/get/", handlers.GetUserHandler(userService))
	mux.HandleFunc("/update/", handlers.UpdateUserHandler(userService))
	mux.HandleFunc("/delete/", handlers.DeleteUserHandler(userService))
	mux.HandleFunc("/users", handlers.GetAllUserHandler(userService))

	mid := middleware.LoggingMiddleware(logger, mux)

	return mid
}
