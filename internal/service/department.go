package service

import (
	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/repository"
)

type Department interface {
	Create(d *models.Department) (*models.Department, error)
	GetByID(id int64, depth uint, includeEmployees bool, sort, order string) (*models.DepartmentResponse, error)
}

type department struct {
	r repository.Repository
}

func NewDepartmentService(r repository.Repository) Department {
	return &department{r: r}
}

func (s *department) Create(d *models.Department) (*models.Department, error) {
	return s.r.Department.Create(d)
}

func (s *department) GetByID(id int64, depth uint, includeEmployees bool, sort, order string) (*models.DepartmentResponse, error) {
	var defaultSort, defaultOrder string

	if order == "" || (order != "ASC" && order != "DESC") {
		defaultOrder = "DESC"
	} else {
		defaultOrder = order
	}

	if sort == "" || (sort != "created_at" && sort != "full_name") {
		defaultSort = "created_at"
	} else {
		defaultSort = sort
	}
	return s.tree(id, depth, includeEmployees, defaultSort, defaultOrder)
}

func (s *department) tree(id int64, depth uint, includeEmployees bool, sort, order string) (*models.DepartmentResponse, error) {
	depart, err := s.r.Department.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &models.DepartmentResponse{
		Department: *depart,
		Children:   []models.DepartmentResponse{},
	}

	if includeEmployees {
		employees, err := s.r.Employee.GetByDepartmentID(id, sort+" "+order)
		if err != nil {
			return nil, err
		}
		response.Employees = employees
	}

	if depth > 0 {
		children, err := s.r.Department.GetChildren(id)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			childResp, err := s.tree(child.ID, depth-1, includeEmployees, sort, order)
			if err != nil {
				return nil, err
			}
			response.Children = append(response.Children, *childResp)
		}
	}

	return response, nil
}
