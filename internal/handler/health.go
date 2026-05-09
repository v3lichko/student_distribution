package handler

import (
	"net/http"

	"github.com/v3lichko/student-distribution/internal/response"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
