package service

import (
	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/repository"
)

type employee struct {
	r repository.Employee
}

type Employee interface {
	Create(e *models.Employee) (*models.Employee, error)
}

func NewEmployeeService(r repository.Employee) Employee {
	return &employee{r: r}
}

func (s *employee) Create(e *models.Employee) (*models.Employee, error) {
	return s.r.Create(e)
}
