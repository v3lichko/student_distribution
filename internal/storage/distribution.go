package storage

import (
	"sort"

	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
)

type DistributionStorage struct {
	db *pg.DB
}

func NewDistributionStorage(db *pg.DB) *DistributionStorage {
	return &DistributionStorage{
		db: db,
	}
}

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

func (s *DistributionStorage) RunDistribution() []models.Student {
	students := make([]models.Student, 0)
	groups := make([]models.Group, 0)

	s.db.Model(&students).Select()
	s.db.Model(&groups).Select()

	sortByScore(students)
	sortByGroup(groups)

	studentIndex := 0

	for _, group := range groups {
		for idx := 0; idx < group.Capacity && studentIndex < len(students); idx++ {
			students[studentIndex].GroupNumber = &group.Number

			s.db.Model(&students[studentIndex]).
				Column("group_number").
				Where("isu = ?", students[studentIndex].ISU).
				Update()

			studentIndex++
		}
	}

	return students
}

func (s *DistributionStorage) GetDistribution() []models.GroupDistribution {
	students := make([]models.Student, 0)

	s.db.Model(&students).
		Where("group_number IS NOT NULL").
		Order("group_number ASC").
		Order("score DESC").
		Select()

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

	return result
}

func (s *DistributionStorage) GetDistributedStudents() []models.Student {
	students := make([]models.Student, 0)

	s.db.Model(&students).
		Where("group_number IS NOT NULL").
		Order("group_number ASC").
		Order("score DESC").
		Select()

	return students
}
