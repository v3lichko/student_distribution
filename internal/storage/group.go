package storage

import (
	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/models"
)

type GroupStorage struct {
	db *pg.DB
}

func NewGroupStorage(db *pg.DB) *GroupStorage {
	return &GroupStorage{
		db: db,
	}
}

func (s *GroupStorage) CreateGroup(group *models.Group) error {
	_, err := s.db.Model(group).Insert()
	return err
}

func (s *GroupStorage) GetGroups() ([]models.Group, error) {
	groups := make([]models.Group, 0)
	err := s.db.Model(&groups).Select()
	return groups, err
}
