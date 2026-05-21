package repository

import "gorm.io/gorm"

type Repository struct {
	Department Department
	Employee   Employee
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		Department: NewDepartmentRepository(db),
		Employee:   NewEmployeeRepository(db),
	}
}
