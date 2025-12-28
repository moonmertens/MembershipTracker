package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func addMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var m Member

	// Attempt to get member from JSON data
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Check for valid phone number
	if !validatePhoneNumber(m.PhoneNumber) {
		respondError(w, http.StatusBadRequest, "Invalid phone number. Must be 8 digits and start with 8 or 9.")
		return
	}

	// Check for valid name
	if m.Name == "" {
		respondError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	insertSQL := `INSERT INTO members (phone_number, name, visits) VALUES ($1, $2, $3)`
	_, err := db.Exec(insertSQL, m.PhoneNumber, m.Name, m.Visits)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to add member (might already exist)")
		return
	}

	msg := fmt.Sprintf("Name: %s\nVisits: %d", m.Name, m.Visits)
	// Prepend 65 to 8-digit number
	fullPhone := m.PhoneNumber + 6500000000
	go SendWhatsApp(fullPhone, msg, nil)

	respondJSON(w, http.StatusCreated, map[string]string{"message": "Member added successfully"})
}

func getMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	phoneStr := r.URL.Query().Get("phone_number")
	if phoneStr == "" {
		respondError(w, http.StatusBadRequest, "Missing 'phone_number' parameter")
		return
	}

	phoneNumber, err := strconv.Atoi(phoneStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	var m Member
	query := `SELECT phone_number, name, visits FROM members WHERE phone_number = $1`
	err = db.QueryRow(query, phoneNumber).Scan(&m.PhoneNumber, &m.Name, &m.Visits)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Member not found")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	respondJSON(w, http.StatusOK, m)
}

func updateMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Only PUT method is allowed")
		return
	}

	var m Member
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	if m.Name == "" {
		respondError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	updateSQL := `UPDATE members SET name = $1, visits = $2 WHERE phone_number = $3`
	result, err := db.Exec(updateSQL, m.Name, m.Visits, m.PhoneNumber)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Member not found")
		return
	}

	msg := fmt.Sprintf("Name: %s\nVisits: %d", m.Name, m.Visits)
	// Prepend 65 to 8-digit number
	fullPhone := m.PhoneNumber + 6500000000
	go SendWhatsApp(fullPhone, msg, nil)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Member updated successfully"})
}

func getAllMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	rows, err := db.Query("SELECT phone_number, name, visits FROM members ORDER BY phone_number ASC")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.PhoneNumber, &m.Name, &m.Visits); err != nil {
			respondError(w, http.StatusInternalServerError, "Error scanning database result")
			return
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		respondError(w, http.StatusInternalServerError, "Database iteration error")
		return
	}

	// Return empty list instead of null if no members
	if members == nil {
		members = []Member{}
	}

	respondJSON(w, http.StatusOK, members)
}

func deleteMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Only DELETE method is allowed")
		return
	}

	phoneStr := r.URL.Query().Get("phone_number")
	if phoneStr == "" {
		respondError(w, http.StatusBadRequest, "Missing 'phone_number' parameter")
		return
	}

	phoneNumber, err := strconv.Atoi(phoneStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	deleteSQL := `DELETE FROM members WHERE phone_number = $1`
	_, err = db.Exec(deleteSQL, phoneNumber)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	msg := "You have been deleted from our members list"
	// Prepend 65 to 8-digit number
	fullPhone := phoneNumber + 6500000000
	go SendWhatsApp(fullPhone, msg, nil)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Member deleted successfully"})
}

// Debugging
// func sendTestMessage(w http.ResponseWriter, r *http.Request) {
//     // 1. Get the phone number from the URL query (e.g., ?phone=1234567890)
//     phoneStr := r.URL.Query().Get("phone")
//     if phoneStr == "" {
//         http.Error(w, "Missing 'phone' parameter", http.StatusBadRequest)
//         return
//     }

//     // 2. Convert to integer
//     phone, err := strconv.Atoi(phoneStr)
//     if err != nil {
//         http.Error(w, "Invalid phone number", http.StatusBadRequest)
//         return
//     }

//     // 3. Send the message
//     msg := "This is a test message from your Membership Tracker Backend! 🚀"
//     SendWhatsApp(phone, msg)

//     // 4. Respond to browser
//     w.Write([]byte("Test message sent to " + phoneStr))
// }
