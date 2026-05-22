package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"gorm.io/gorm"
)

type createEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	depID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var body createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	employee, err := h.s.Employee.Create(&models.Employee{
		DepartmentID: depID,
		FullName:     body.FullName,
		Position:     body.Position,
		HiredAt:      body.HiredAt,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "department not found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, employee)
}
