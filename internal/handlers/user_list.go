package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ultrabor/go-user-api/internal/services"
)

func GetAllUserHandler(svc *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var (
			limit = 0
			page  = 0
			age   *int
		)
		var name *string
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
			*age, err = strconv.Atoi(query.Get("age"))
			if err != nil || *age <= 0 {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
		}

		if query.Has("name") {
			*name = query.Get("name")
		}

		users, err := svc.GetAllUsers(limit, page, name, age)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
