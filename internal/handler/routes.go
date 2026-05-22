package handler

import "net/http"

func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /departments/", h.CreateDepartment)
	mux.HandleFunc("GET /departments/{id}", h.GetDepartment)
	mux.HandleFunc("PATCH /departments/{id}", h.UpdateDepartment)
	mux.HandleFunc("DELETE /departments/{id}", h.DeleteDepartment)
	mux.HandleFunc("POST /departments/{id}/employees/", h.CreateEmployee)

	return mux
}
