package handler

import (
	"net/http"

	"github.com/v3lichko/student-distribution/internal/response"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
