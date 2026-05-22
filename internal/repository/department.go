package repository

import (
	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"gorm.io/gorm"
)

type Department interface {
	Create(d *models.Department) (*models.Department, error)
	GetByID(id int64) (*models.Department, error)
	GetChildren(parentID int64) ([]models.Department, error)
	Update(id int64, updates map[string]interface{}) (*models.Department, error)
	Delete(id int64) error
	HasChildren(id int64) (bool, error)
	DeleteWithReassign(id, toDeptID int64) error
	GetByName(name string, parentID *int64) (*models.Department, error)
}

type department struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) Department {
	return &department{db: db}
}

func (r *department) Create(d *models.Department) (*models.Department, error) {
	if err := r.db.Create(d).Error; err != nil {
		return nil, err
	}

	return r.GetByID(d.ID)
}

func (r *department) GetByID(id int64) (*models.Department, error) {
	var dept models.Department
	if err := r.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *department) GetChildren(parentID int64) ([]models.Department, error) {
	var children []models.Department
	if err := r.db.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}

func (r *department) Update(id int64, updates map[string]interface{}) (*models.Department, error) {

	if err := r.db.Model(&models.Department{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *department) GetByName(name string, parentID *int64) (*models.Department, error) {
	var d models.Department
	q := r.db.Where("name = ?", name)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}
	if err := q.First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *department) Delete(id int64) error {
	return r.db.Delete(&models.Department{}, id).Error
}

func (r *department) HasChildren(id int64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Department{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *department) DeleteWithReassign(id, toDeptID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Employee{}).Where("department_id = ?", id).Update("department_id", toDeptID).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Department{}, id).Error
	})
}
