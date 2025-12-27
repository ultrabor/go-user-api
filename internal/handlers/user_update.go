package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	m "github.com/ultrabor/go-user-api/internal/models"
	"github.com/ultrabor/go-user-api/internal/services"
)

func UpdateUserHandler(svc *services.UserService) http.HandlerFunc {
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
			Name *string `json:"name"`
			Age  *int    `json:"age"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		user := m.User{ID: id, Name: "", Age: 0}

		if input.Name != nil {
			user.Name = *input.Name
		}
		if input.Age != nil {
			user.Age = *input.Age
		}

		updatedUser, err := svc.UpdateUser(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)
	}
}
