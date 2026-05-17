package distribution

import (
	"sort"

	"github.com/v3lichko/student-distribution/internal/models"
)

func Distribute(students []models.Student, groups []models.Group) []models.Student {
	sort.Slice(students, func(i, j int) bool {
		return students[i].Score > students[j].Score
	})

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Number < groups[j].Number
	})

	studentIndex := 0

	for _, group := range groups {
		for i := 0; i < group.Capacity && studentIndex < len(students); i++ {
			groupNumber := group.Number
			students[studentIndex].GroupNumber = &groupNumber
			studentIndex++
		}
	}

	return students
}
