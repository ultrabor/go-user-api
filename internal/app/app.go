package app

import (
	"log/slog"
	"net/http"
	"os"

	handlers "github.com/ultrabor/go-user-api/internal/handlers"
	mw "github.com/ultrabor/go-user-api/internal/middleware"
	st "github.com/ultrabor/go-user-api/internal/storage"
)

func RunApp() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	store := st.NewStore()

	mux := http.NewServeMux()

	mux.HandleFunc("/create", handlers.CreateUserHandler(store))
	mux.HandleFunc("/get/", handlers.GetUserHandler(store))
	mux.HandleFunc("/update", handlers.UpdateUserHandler(store))
	mux.HandleFunc("/delete/", handlers.DeleteUserHandler(store))

	mid := mw.LoggingMiddleware(logger, mux)

	http.ListenAndServe(":8080", mid)
}
