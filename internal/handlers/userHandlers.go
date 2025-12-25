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

		idStr := strings.TrimPrefix(r.URL.Path, "/get/")

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

func GetAllUserHandler(st *s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var (
			limit = 0
			page  = 0
			age   = -1
		)
		var name string
		var err error

		query := r.URL.Query()

		if query.Has("limit") {
			limit, err = strconv.Atoi(query.Get("limit"))
			if err != nil || limit <= 0 {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
		}

		if query.Has("page") {
			page, err = strconv.Atoi(query.Get("page"))
			if err != nil || page <= 0 {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
		}

		if query.Has("age") {
			age, err = strconv.Atoi(query.Get("age"))
			if err != nil || age <= 0 {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
		}

		if query.Has("name") {
			name = query.Get("name")
		}

		users := st.GetAllUser(limit, page, name, age)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func DeleteUserHandler(st *s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/delete/")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user, ok := st.DeleteUser(id)

		if !ok {
			http.Error(w, "bad request", http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	}
}

func UpdateUserHandler(st *s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/update/")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "invalid user id", http.StatusBadRequest)
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

		if input.Name == "" || input.Age < 0 {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		user, ok := st.UpdateUser(id, input.Name, input.Age)

		if !ok {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
