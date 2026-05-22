package service

import (
	"errors"
	"strings"

	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrSelfParent    = errors.New("department cannot be its own parent")
	ErrOwnDescendant = errors.New("cannot move department into its own subtree")
	ErrNameConflict  = errors.New("name already exists in this parent")
)

type Department interface {
	Create(d *models.Department) (*models.Department, error)
	GetByID(id int64, depth uint, includeEmployees bool, sort, order string) (*models.DepartmentResponse, error)
	Update(id int64, dto *models.UpdateDepartment) (*models.Department, error)
	Delete(id int64, mode string, reassignTo *int64) error
}

type department struct {
	r repository.Repository
}

func NewDepartmentService(r repository.Repository) Department {
	return &department{r: r}
}

func (s *department) isDescendant(deptID, parentID int64) (bool, error) {
	curr := parentID
	for {
		dept, err := s.r.Department.GetByID(curr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		if dept.ParentID == nil {
			return false, nil
		}
		if *dept.ParentID == deptID {
			return true, nil
		}
		curr = *dept.ParentID
	}
}

func (s *department) Create(d *models.Department) (*models.Department, error) {
	d.Name = strings.TrimSpace(d.Name)
	if strings.TrimSpace(d.Name) == "" {
		return nil, errors.New("name required")
	}
	if len([]rune(d.Name)) > 200 || len([]rune(d.Name)) < 1 {
		return nil, errors.New("name max 200")
	}
	_, err := s.r.Department.GetByName(d.Name, d.ParentID)
	if err == nil {
		return nil, ErrNameConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

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

func (s *department) Delete(id int64, mode string, to *int64) error {
	if _, err := s.r.Department.GetByID(id); err != nil {
		return err
	}

	switch mode {
	case "cascade":
		return s.r.Department.Delete(id)
	case "reassign":
		if to == nil {
			return errors.New("reassign_to_department_id required")
		}
		if _, err := s.r.Department.GetByID(*to); err != nil {
			return err
		}
		hasChildren, err := s.r.Department.HasChildren(id)
		if err != nil {
			return err
		}
		if hasChildren {
			return errors.New("department has children")
		}
		return s.r.Department.DeleteWithReassign(id, *to)
	default:
		return errors.New("cascade or reassign")
	}
}

func (s *department) Update(id int64, dto *models.UpdateDepartment) (*models.Department, error) {
	curr, err := s.r.Department.GetByID(id)
	if err != nil {
		return nil, err
	}

	upd := make(map[string]interface{})

	if dto.Name != nil {
		name := strings.TrimSpace(*dto.Name)
		if name == "" {
			return nil, errors.New("name required")
		}
		if len([]rune(name)) > 200 {
			return nil, errors.New("name max 200")
		}
		parentID := curr.ParentID
		if dto.ParentID != nil {
			parentID = dto.ParentID
		}
		existing, err := s.r.Department.GetByName(name, parentID)
		if err == nil && existing.ID != id {
			return nil, ErrNameConflict
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		upd["name"] = name
	}

	if dto.ParentID != nil {
		if id == *dto.ParentID {
			return nil, ErrSelfParent
		}
		isDesc, err := s.isDescendant(id, *dto.ParentID)
		if err != nil {
			return nil, err
		}
		if isDesc {
			return nil, ErrOwnDescendant
		}
		upd["parent_id"] = *dto.ParentID
	}

	if len(upd) == 0 {
		return curr, nil
	}

	return s.r.Department.Update(id, upd)
}
