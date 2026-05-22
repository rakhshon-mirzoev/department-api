package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/service"
	"gorm.io/gorm"
)

type createDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var body createDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	dept, err := h.s.Department.Create(&models.Department{
		Name:     body.Name,
		ParentID: body.ParentID,
	})
	if err != nil {
		if errors.Is(err, service.ErrNameConflict) {
			writeError(w, http.StatusConflict, err.Error())
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, dept)
}

func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	var (
		depth            uint
		includeEmployees bool
		depID            int64
		sort, order      string
	)
	sort = r.URL.Query().Get("sort")
	order = r.URL.Query().Get("order")

	depIDStr := r.PathValue("id")
	depID, err := strconv.ParseInt(depIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	depthStr := r.URL.Query().Get("depth")
	if depthStr == "" {
		depth = 1
	} else {
		depthInt, err := strconv.ParseUint(depthStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid depth")
			return
		}
		depth = uint(depthInt)
	}
	if depth > 5 {
		writeError(w, http.StatusBadRequest, "depth max 5")
		return
	}

	includeEmployeesStr := r.URL.Query().Get("include_employees")
	if includeEmployeesStr == "" {
		includeEmployees = true
	} else {
		includeEmployeesBool, err := strconv.ParseBool(includeEmployeesStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid include_employees")
			return
		}
		includeEmployees = includeEmployeesBool
	}

	departmentWithChildrenNEmployees, err := h.s.Department.GetByID(depID, depth, includeEmployees, sort, order)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "department not found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, departmentWithChildrenNEmployees)
}

func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var data models.UpdateDepartment
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	updatedDepart, err := h.s.Department.Update(id, &data)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			writeError(w, http.StatusNotFound, "department not found")
		case errors.Is(err, service.ErrOwnDescendant),
			errors.Is(err, service.ErrSelfParent),
			errors.Is(err, service.ErrNameConflict):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, updatedDepart)
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	mode := r.URL.Query().Get("mode")
	switch mode {
	case "cascade":
		if err := h.s.Department.Delete(id, mode, nil); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeError(w, http.StatusNotFound, "department not found")
			} else {
				writeError(w, http.StatusBadRequest, err.Error())
			}
			return
		}
	case "reassign":
		reassignTo, err := strconv.ParseInt(r.URL.Query().Get("reassign_to_department_id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid reassign id")
			return
		}
		if err := h.s.Department.Delete(id, mode, &reassignTo); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				writeError(w, http.StatusNotFound, "department not found")
			} else {
				writeError(w, http.StatusBadRequest, err.Error())
			}
			return
		}
	default:
		writeError(w, http.StatusBadRequest, "mode must be cascade or reassign")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
