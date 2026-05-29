package api

import "github.com/v3lichko/student-distribution/internal/models"

func StudentFromModel(student models.Student) Student {
	return Student{
		ISU:         student.ISU,
		FullName:    student.FullName,
		Telegram:    student.Telegram,
		Score:       student.Score,
		GroupNumber: student.GroupNumber,
	}
}

func StudentsFromModels(students []models.Student) []Student {
	result := make([]Student, 0, len(students))

	for _, student := range students {
		result = append(result, StudentFromModel(student))
	}

	return result
}
