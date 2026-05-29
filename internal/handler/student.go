package handler

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/v3lichko/student-distribution/internal/api"
	"github.com/v3lichko/student-distribution/internal/models"
	"github.com/v3lichko/student-distribution/internal/response"
	"github.com/v3lichko/student-distribution/internal/storage"
)

type StudentHandler struct {
	storage *storage.StudentStorage
}

func NewStudentHandler(studentStorage *storage.StudentStorage) *StudentHandler {
	return &StudentHandler{storage: studentStorage}
}

// @Summary Create student
// @Tags students
// @Accept json
// @Produce json
// @Param body body models.Student true "Student data"
// @Success 201 {object} api.Student
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students [post]
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if err := h.storage.CreateStudent(&student); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, api.StudentFromModel(student))
}

// @Summary Get all students
// @Tags students
// @Produce json
// @Success 200 {array} api.Student
// @Failure 500 {object} map[string]string
// @Router /students [get]
func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.storage.GetStudents()
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	response.WriteJSON(w, http.StatusOK, api.StudentsFromModels(students))
}

// @Summary Delete student
// @Tags students
// @Produce json
// @Param isu query int true "ISU number"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students [delete]
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	isu, err := strconv.Atoi(r.URL.Query().Get("isu"))
	if err != nil || isu <= 0 {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid isu"})
		return
	}

	if err := h.storage.DeleteStudent(isu); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// @Summary Import students from CSV
// @Tags students
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /students/import [post]
func (h *StudentHandler) ImportStudentsCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "file required"})
		return
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid csv"})
		return
	}

	for i, record := range records {
		if i == 0 || len(record) < 4 {
			continue
		}

		isu, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		score, err := strconv.Atoi(record[3])
		if err != nil {
			continue
		}

		_ = h.storage.CreateStudent(&models.Student{
			ISU:      isu,
			FullName: record[1],
			Telegram: record[2],
			Score:    score,
		})
	}

	response.WriteJSON(w, http.StatusCreated, map[string]string{"status": "students imported"})
}
