package repository

import (
	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"gorm.io/gorm"
)

type Employee interface {
	Create(e *models.Employee) error
}

type employee struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) Employee {
	return &employee{db: db}
}

func (r *employee) Create(e *models.Employee) error {
	if err := r.db.Create(e).Error; err != nil {
		return err
	}
	return nil
}
