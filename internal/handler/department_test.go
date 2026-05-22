package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rakhshon-mirzoev/department-api/internal/models"
	"github.com/rakhshon-mirzoev/department-api/internal/service"
)

type mockDepartmentService struct {
	createFn func(d *models.Department) (*models.Department, error)
}

func (m *mockDepartmentService) Create(d *models.Department) (*models.Department, error) {
	return m.createFn(d)
}
func (m *mockDepartmentService) GetByID(id int64, depth uint, includeEmployees bool, sort, order string) (*models.DepartmentResponse, error) {
	return nil, nil
}
func (m *mockDepartmentService) Update(id int64, dto *models.UpdateDepartment) (*models.Department, error) {
	return nil, nil
}
func (m *mockDepartmentService) Delete(id int64, mode string, reassignTo *int64) error {
	return nil
}

func TestCreateDepartment(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		mockResult *models.Department
		mockErr    error
		wantStatus int
	}{
		{
			name:       "успешное создание",
			body:       map[string]any{"name": "Backend"},
			mockResult: &models.Department{ID: 1, Name: "Backend"},
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "пустое имя",
			body:       map[string]any{"name": ""},
			mockResult: nil,
			mockErr:    errors.New("name required"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "конфликт имени",
			body:       map[string]any{"name": "Backend"},
			mockResult: nil,
			mockErr:    service.ErrNameConflict,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "невалидный json",
			body:       "это не json",
			mockResult: nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bodyBytes, _ := json.Marshal(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/departments/", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			h := &Handler{
				s: service.Service{
					Department: &mockDepartmentService{
						createFn: func(d *models.Department) (*models.Department, error) {
							return tt.mockResult, tt.mockErr
						},
					},
				},
			}

			h.CreateDepartment(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("тест %q: ожидали статус %d, получили %d", tt.name, tt.wantStatus, w.Code)
			}
		})
	}
}
