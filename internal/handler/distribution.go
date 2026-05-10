package handler

import (
	"encoding/csv"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
	"github.com/v3lichko/student-distribution/internal/response"
)

func sortByScore(students []models.Student) {
	sort.Slice(students, func(i int, j int) bool {
		return students[i].Score > students[j].Score
	})
}

func sortByGroup(groups []models.Group) {
	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Number < groups[j].Number
	})
}

func sortDistributionByGroup(result []models.GroupDistribution) {
	sort.Slice(result, func(i int, j int) bool {
		return result[i].GroupNumber < result[j].GroupNumber
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

func (h *DistributionHandler) Distribution(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetDistribution(w, r)
		return
	}
}

func (h *DistributionHandler) GetDistribution(w http.ResponseWriter, r *http.Request) {
	students := make([]models.Student, 0)
	h.db.Model(&students).Where("group_number IS NOT NULL").Order("group_number ASC").Order("score DESC").Select()
	resultMap := make(map[int][]models.Student)

	for _, student := range students {
		if student.GroupNumber == nil {
			continue
		}
		groupNumber := *student.GroupNumber
		resultMap[groupNumber] = append(resultMap[groupNumber], student)
	}

	result := make([]models.GroupDistribution, 0)
	for groupNumber, groupStudents := range resultMap {
		result = append(result, models.GroupDistribution{
			GroupNumber: groupNumber,
			Students:    groupStudents,
		})
	}
	sortDistributionByGroup(result)
	response.WriteJSON(w, http.StatusOK, result)
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

func (h *DistributionHandler) ExportDistributionCSV(w http.ResponseWriter, r *http.Request) {
	students := make([]models.Student, 0)
	h.db.Model(&students).
		Where("group_number IS NOT NULL").
		Order("group_number ASC").
		Order("score DESC").
		Select()

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
