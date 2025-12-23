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

	mux.HandleFunc("/users", handlers.CreateUserHandler(store))
	mux.HandleFunc("/users/", handlers.GetUserHandler(store))

	mid := mw.LoggingMiddleware(logger, mux)

	http.ListenAndServe(":8080", mid)
}
