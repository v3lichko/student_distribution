package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"

	distribution "github.com/v3lichko/student-distribution/internal/distributition"
	"github.com/v3lichko/student-distribution/internal/response"
	"github.com/v3lichko/student-distribution/internal/storage"
)

type DistributionHandler struct {
	storage *storage.DistributionStorage
}

func NewDistributionHandler(distributionStorage *storage.DistributionStorage) *DistributionHandler {
	return &DistributionHandler{storage: distributionStorage}
}

// @Summary Get current distribution
// @Tags distribution
// @Produce json
// @Success 200 {array} models.GroupDistribution
// @Failure 500 {object} map[string]string
// @Router /distribution [get]
func (h *DistributionHandler) GetDistribution(w http.ResponseWriter, r *http.Request) {
	result, err := h.storage.GetDistribution(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

// @Summary Run distribution algorithm
// @Tags distribution
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {object} map[string]string
// @Router /distribution/run [post]
func (h *DistributionHandler) StartDistribution(w http.ResponseWriter, r *http.Request) {
	students, err := h.storage.ListStudents(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	groups, err := h.storage.ListGroups(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	assigned := distribution.Distribute(students, groups)

	if err := h.storage.UpdateAssignments(r.Context(), assigned); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, assigned)
}

// @Summary Export distribution as CSV
// @Tags distribution
// @Produce text/csv
// @Success 200
// @Failure 500 {object} map[string]string
// @Router /distribution/export [get]
func (h *DistributionHandler) ExportDistributionCSV(w http.ResponseWriter, r *http.Request) {
	students, err := h.storage.GetDistributedStudents(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="distribution.csv"`)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	_ = writer.Write([]string{"group_number", "isu", "full_name", "telegram", "score"})

	for _, student := range students {
		groupNumber := ""
		if student.GroupNumber != nil {
			groupNumber = strconv.Itoa(*student.GroupNumber)
		}

		_ = writer.Write([]string{
			groupNumber,
			strconv.Itoa(student.ISU),
			student.FullName,
			student.Telegram,
			strconv.Itoa(student.Score),
		})
	}
}
