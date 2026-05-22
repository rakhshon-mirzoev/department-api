package service

import "github.com/rakhshon-mirzoev/department-api/internal/repository"

type Service struct {
	Employee   Employee
	Department Department
}

func NewService(r repository.Repository) Service {
	return Service{
		Employee:   NewEmployeeService(r),
		Department: NewDepartmentService(r),
	}
}
