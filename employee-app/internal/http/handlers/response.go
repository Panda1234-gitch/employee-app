package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"unicode"
)

// common JSON response helper
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// username cannot be only numbers
func isNumericOnly(s string) bool {
	if s == "" {
		return false
	}
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

// password checks:
// - min 6
// - at least 1 letter, 1 digit, 1 special
func validatePassword(pw string) (bool, string) {
	if len(pw) < 6 {
		return false, "Password must be at least 6 characters long"
	}

	hasLetter := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range pw {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsDigit(ch):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if !hasLetter {
		return false, "Password must contain at least one letter"
	}
	if !hasDigit {
		return false, "Password must contain at least one digit"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special symbol"
	}
	return true, ""
}
