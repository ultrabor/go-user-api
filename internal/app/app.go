package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ultrabor/go-user-api/internal/config"
	"github.com/ultrabor/go-user-api/internal/server"

	storage "github.com/ultrabor/go-user-api/internal/storage"
	memory "github.com/ultrabor/go-user-api/internal/storage/memory"
	postgres "github.com/ultrabor/go-user-api/internal/storage/postgres"

	_ "github.com/lib/pq"
)

func RunApp() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	usePostgres := true

	var store storage.UserStore

	if usePostgres {
		s, err := postgres.New(config.GetPostgresDSN(), logger)
		if err != nil {
			panic(err)
		}
		store = s
	} else {
		store = memory.New(logger)
	}

	ser := server.Server(logger, store)

	http.ListenAndServe(":8080", ser)
}
