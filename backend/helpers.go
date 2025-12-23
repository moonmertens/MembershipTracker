package main

import (
    "encoding/json"
    "net/http"
)

// respondJSON sends a JSON response with a specific status code
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

// respondError sends a JSON error message
func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}

// Function to validate phone number in Singaporean context
func validatePhoneNumber(phone int) bool {
    return phone >= 80000000 && phone <= 99999999
}