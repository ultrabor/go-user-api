package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	s "github.com/ultrabor/go-user-api/internal/storage"
)

func CreateUserHandler(st *s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user := st.CreateUser(input.Name, input.Age)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func GetUserHandler(st *s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/users/")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, ok := st.GetUser(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
