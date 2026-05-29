package storage

import (
	"context"
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

func (s *DistributionStorage) ListStudents(ctx context.Context) ([]models.Student, error) {
	students := make([]models.Student, 0)

	err := s.db.ModelContext(ctx, &students).Select()
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (s *DistributionStorage) ListGroups(ctx context.Context) ([]models.Group, error) {
	groups := make([]models.Group, 0)

	err := s.db.ModelContext(ctx, &groups).Select()
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *DistributionStorage) UpdateAssignments(ctx context.Context, students []models.Student) error {
	for i := range students {
		_, err := s.db.ModelContext(ctx, &students[i]).
			Column("group_number").
			Where("isu = ?", students[i].ISU).
			Update()

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DistributionStorage) GetDistribution(ctx context.Context) ([]models.GroupDistribution, error) {
	students := make([]models.Student, 0)

	err := s.db.ModelContext(ctx, &students).
		Where("group_number IS NOT NULL").
		Order("group_number ASC").
		Order("score DESC").
		Select()

	if err != nil {
		return nil, err
	}

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

	sort.Slice(result, func(i, j int) bool {
		return result[i].GroupNumber < result[j].GroupNumber
	})

	return result, nil
}

func (s *DistributionStorage) GetDistributedStudents(ctx context.Context) ([]models.Student, error) {
	students := make([]models.Student, 0)

	err := s.db.ModelContext(ctx, &students).
		Where("group_number IS NOT NULL").
		Order("group_number ASC").
		Order("score DESC").
		Select()

	if err != nil {
		return nil, err
	}

	return students, nil
}
