package repository

import (
	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"gorm.io/gorm"
)

type Department interface {
	Create(d *models.Department) error
}

type department struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) Department {
	return &department{db: db}
}

func (r *department) Create(d *models.Department) error {
	if err := r.db.Create(d).Error; err != nil {
		return err
	}
	return nil
}
