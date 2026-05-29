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

func (s *StudentStorage) CreateStudent(student *models.Student) error {
	_, err := s.db.Model(student).Insert()
	return err
}

func (s *StudentStorage) GetStudents() ([]models.Student, error) {
	students := make([]models.Student, 0)
	err := s.db.Model(&students).Select()
	return students, err
}

func (s *StudentStorage) DeleteStudent(isu int) error {
	_, err := s.db.Model((*models.Student)(nil)).
		Where("isu = ?", isu).
		Delete()
	return err
}
