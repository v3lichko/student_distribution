package handler

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *StudentHandler) Students(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateStudent(w, r)
		return
	}
	if r.Method == http.MethodGet {
		h.GetStudents(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		h.DeleteStudent(w, r)
		return
	}

	response.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error": "method not allowed",
	})
}

// @Summary Delete students
// @Tags students
// @Produce json
// @Param isu query int true "isu of student"
// @Success 200 {object} map[string]string
// @Router /students [delete]
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	isuStr := r.URL.Query().Get("isu")
	isu, _ := strconv.Atoi(isuStr)
	_, err := h.db.Model((*models.Student)(nil)).
		Where("isu = ?", isu).
		Delete()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
}

// @Summary Get students
// @Tags students
// @Produce json
// @Success 200 {array} map[string]string
// @Router /students [get]
func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	student := make([]models.Student, 0)
	h.db.Model(&student).Select()
	response.WriteJSON(w, http.StatusOK, student)
}

// @Summary Import students
// @Tags students
// @Accept multipart/form-data
// @Produce json
// @Param film formData file true "CSV with students"
// @Success 201 {object} map[string]string
// @Router /students/import [post]
func (h *StudentHandler) ImportStudentsCSV(w http.ResponseWriter, r *http.Request) {
	file, _, _ := r.FormFile("file")
	defer file.Close()
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	for i, record := range records {
		if i == 0 {
			continue
		}

		isu, _ := strconv.Atoi(record[0])
		score, _ := strconv.Atoi(record[3])

		student := models.Student{
			ISU:      isu,
			FullName: record[1],
			Telegram: record[2],
			Score:    score,
		}
		h.db.Model(&student).Insert()
	}
	response.WriteJSON(w, http.StatusCreated, map[string]string{
		"status": "students imported",
	})
}
