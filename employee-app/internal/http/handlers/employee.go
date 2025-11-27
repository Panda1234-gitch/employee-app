package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"employee-app/internal/repository"
)

type EmployeeHandler struct {
	Repo repository.EmployeeRepository
}

func NewEmployeeHandler(repo repository.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{Repo: repo}
}


// GetByID godoc
// @Summary Get employee by ID
// @Tags Employee
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} map[string]any
// @Failure 404 {object} map[string]string
// @Router /employees/{id} [get]
func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Only GET allowed"})
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "employees" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL"})
		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	emp, err := h.Repo.GetByID(id)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "DB query failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"employee": emp})
}
