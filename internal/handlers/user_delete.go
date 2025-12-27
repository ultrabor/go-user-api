package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ultrabor/go-user-api/internal/services"
)

func DeleteUserHandler(svc *services.UserService) http.HandlerFunc {
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

		err = svc.DeleteUser(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
