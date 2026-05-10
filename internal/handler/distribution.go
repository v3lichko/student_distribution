package handler

import (
	"net/http"
	"sort"

	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
	"github.com/v3lichko/student-distribution/internal/response"
)

func sortByScore(students []models.Student) {
	sort.Slice(students, func(i int, j int) bool {
		return students[i].Score > students[j].Score
	})
}

func sortByGroup(group []models.Group) {
	sort.Slice(group, func(i int, j int) bool {
		return group[i].Number < group[j].Number
	})
}

type DistributionHandler struct {
	db *pg.DB
}

func NewDistributionHandler(db *pg.DB) *DistributionHandler {
	return &DistributionHandler{
		db: db,
	}
}

func (h *DistributionHandler) StartDistribution(w http.ResponseWriter, r *http.Request) {
	students := make([]models.Student, 0)
	groups := make([]models.Group, 0)
	h.db.Model(&students).Select()
	h.db.Model(&groups).Select()
	sortByScore(students)
	sortByGroup(groups)
	studentIndex := 0
	for _, group := range groups {
		for idx := 0; idx < group.Capacity && studentIndex < len(students); idx++ {
			students[studentIndex].GroupNumber = &group.Number
			h.db.Model(&students[studentIndex]).Column("group_number").Where("isu = ?", students[studentIndex].ISU).Update()
			studentIndex++
		}
	}
	response.WriteJSON(w, http.StatusOK, students)
}
