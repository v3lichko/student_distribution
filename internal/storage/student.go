package storage

import (
	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
)

type StudentStorage struct {
	db *pg.DB
}

func NewStudentStorage(db *pg.DB) *StudentStorage {
	return &StudentStorage{
		db: db,
	}
}

func (s *StudentStorage) CreateStudent(student *models.Student) {
	s.db.Model(student).Insert()
}

func (s *StudentStorage) GetStudents() []models.Student {
	students := make([]models.Student, 0)
	s.db.Model(&students).Select()
	return students
}

func (s *StudentStorage) DeleteStudent(isu int) {
	s.db.Model((*models.Student)(nil)).
		Where("isu = ?", isu).
		Delete()
}
