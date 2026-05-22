package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rakhshon-mirzoev/department-api/internal/service"
)

type Handler struct {
	s service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{s: s}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
