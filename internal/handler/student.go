package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
	"github.com/v3lichko/student-distribution/internal/response"
)

type StudentHandler struct {
	db *pg.DB
}

func NewStudentHandler(db *pg.DB) *StudentHandler {
	return &StudentHandler{
		db: db,
	}
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	json.NewDecoder(r.Body).Decode(&student)

	h.db.Model(&student).Insert()
	response.WriteJSON(w, http.StatusCreated, student)
}
