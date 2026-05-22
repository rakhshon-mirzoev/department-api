package repository

import (
	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"gorm.io/gorm"
)

type Employee interface {
	Create(e *models.Employee) (*models.Employee, error)
	GetByID(id int64) (*models.Employee, error)
	GetByDepartmentID(departmentID int64, order string) ([]models.Employee, error)
}

type employee struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) Employee {
	return &employee{db: db}
}

func (r *employee) Create(e *models.Employee) (*models.Employee, error) {
	if err := r.db.Create(e).Error; err != nil {
		return nil, err
	}
	return r.GetByID(e.ID)
}

func (r *employee) GetByID(id int64) (*models.Employee, error) {
	var emp models.Employee
	if err := r.db.First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employee) GetByDepartmentID(departmentID int64, order string) ([]models.Employee, error) {
	var employees []models.Employee
	if err := r.db.Where("department_id = ?", departmentID).Order(order).Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}
