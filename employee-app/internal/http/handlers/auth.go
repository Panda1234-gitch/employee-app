package handlers

// NOTE: imports unchanged...

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// "employee-app/internal/models"
	"employee-app/internal/repository"
)

type AuthHandler struct {
	Repo repository.EmployeeRepository
}

func NewAuthHandler(repo repository.EmployeeRepository) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

type registerRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Designation string `json:"designation"`
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register a new employee
// @Description Create a new employee
// @Tags Auth
// @Accept json
// @Produce json
// @Param employee body registerRequest true "Employee Data"
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Only POST allowed"})
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Basic required fields
	if req.Name == "" || req.Password == "" || req.Designation == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "All fields required"})
		return
	}

	// Username not only numbers
	if isNumericOnly(req.Name) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Username cannot be only numbers"})
		return
	}

	// Password strength
	if ok, msg := validatePassword(req.Password); !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
		return
	}

	emp, err := h.Repo.Create(req.Name, req.Password, req.Designation)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "DB insert failed"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"employee": emp})
}

// Login godoc
// @Summary Login employee
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body loginRequest true "Login Data"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
// POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Only POST allowed"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	emp, err := h.Repo.GetByCredentials(req.Name, req.Password)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "DB query failed"})
		return
	}

	emp.Password = ""

	writeJSON(w, http.StatusOK, map[string]any{
		"message":  "Login successful",
		"employee": emp,
	})
}
