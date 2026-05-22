package service

import (
	"errors"
	"strings"

	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/repository"
)

type employee struct {
	r repository.Repository
}

type Employee interface {
	Create(e *models.Employee) (*models.Employee, error)
}

func NewEmployeeService(r repository.Repository) Employee {
	return &employee{r: r}
}

func (s *employee) Create(e *models.Employee) (*models.Employee, error) {
	_, err := s.r.Department.GetByID(e.DepartmentID)
	if err != nil {
		return nil, err
	}

	e.FullName = strings.TrimSpace(e.FullName)
	e.Position = strings.TrimSpace(e.Position)
	if err := validateEmployee(e); err != nil {
		return nil, err
	}
	return s.r.Employee.Create(e)
}

func validateEmployee(e *models.Employee) error {
	if strings.TrimSpace(e.FullName) == "" {
		return errors.New("full_name required")
	}
	if len([]rune(e.FullName)) > 200 {
		return errors.New("full_name max 200")
	}
	if strings.TrimSpace(e.Position) == "" {
		return errors.New("position required")
	}
	if len([]rune(e.Position)) > 200 {
		return errors.New("position max 200")
	}
	return nil
}
