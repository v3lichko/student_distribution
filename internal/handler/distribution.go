package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/v3lichko/student-distribution/internal/response"
	"github.com/v3lichko/student-distribution/internal/storage"
)

type DistributionHandler struct {
	storage *storage.DistributionStorage
}

func NewDistributionHandler(distributionStorage *storage.DistributionStorage) *DistributionHandler {
	return &DistributionHandler{
		storage: distributionStorage,
	}
}

func (h *DistributionHandler) Distribution(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetDistribution(w, r)
		return
	}

	response.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error": "method not allowed",
	})
}

func (h *DistributionHandler) GetDistribution(w http.ResponseWriter, r *http.Request) {
	result := h.storage.GetDistribution()

	response.WriteJSON(w, http.StatusOK, result)
}

func (h *DistributionHandler) StartDistribution(w http.ResponseWriter, r *http.Request) {
	students := h.storage.RunDistribution()

	response.WriteJSON(w, http.StatusOK, students)
}

func (h *DistributionHandler) ExportDistributionCSV(w http.ResponseWriter, r *http.Request) {
	students := h.storage.GetDistributedStudents()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="distribution.csv"`)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"group_number",
		"isu",
		"full_name",
		"telegram",
		"score",
	})

	for _, student := range students {
		groupNumber := ""

		if student.GroupNumber != nil {
			groupNumber = strconv.Itoa(*student.GroupNumber)
		}

		writer.Write([]string{
			groupNumber,
			strconv.Itoa(student.ISU),
			student.FullName,
			student.Telegram,
			strconv.Itoa(student.Score),
		})
	}
}
