package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/service"
	"gorm.io/gorm"
)

type mockEmployeeService struct {
	createFn func(e *models.Employee) (*models.Employee, error)
}

func (m *mockEmployeeService) Create(e *models.Employee) (*models.Employee, error) {
	return m.createFn(e)
}

func TestCreateEmployee(t *testing.T) {
	tests := []struct {
		name       string
		deptID     string
		body       any
		mockResult *models.Employee
		mockErr    error
		wantStatus int
	}{
		{
			name:   "успешное создание",
			deptID: "1",
			body: map[string]any{
				"full_name": "Иванов Иван",
				"position":  "developer",
				"hired_at":  "2026-05-22T20:27:10.815147Z",
			},
			mockResult: &models.Employee{ID: 1, FullName: "Иванов Иван", Position: "developer"},
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:   "несуществующий департамент",
			deptID: "9999",
			body: map[string]any{
				"full_name": "Большой Джо",
				"position":  "developer",
			},
			mockResult: nil,
			mockErr:    gorm.ErrRecordNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "невалидный json",
			deptID:     "1",
			body:       "это не json",
			mockResult: nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/departments/"+tt.deptID+"/employees/", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tt.deptID)

			w := httptest.NewRecorder()

			h := &Handler{
				s: service.Service{
					Employee: &mockEmployeeService{
						createFn: func(e *models.Employee) (*models.Employee, error) {
							return tt.mockResult, tt.mockErr
						},
					},
				},
			}

			h.CreateEmployee(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("тест %q: ожидали статус %d, получили %d", tt.name, tt.wantStatus, w.Code)
			}
		})
	}
}
