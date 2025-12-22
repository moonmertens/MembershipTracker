package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

type Member struct {
	PhoneNumber int		`json:"phone_number"`
	Name 		string 	`json:"name"`
	Visits		int 	`json:"visits"`
}

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("sqlite", "./members.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create preliminary table
	// Members are identified by their phone number as primary key
	// Fields are phone number, name, and number of visits
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS members (
		phone_number INTEGER PRIMARY KEY, 
		name TEXT,
		visits INTEGER
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// API endpoints (CRUD)
	http.HandleFunc("/add-member", addMember)
	http.HandleFunc("/get-member", getMember)
	http.HandleFunc("/update-member", updateMember)
	http.HandleFunc("/delete-member", deleteMember)

	fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Function to handle new member creation
func addMember(w http.ResponseWriter, r *http.Request) {
    // Only allow POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var m Member
    // Decode the JSON sent by the user into our 'm' variable
    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }

    // Insert the new member into the database
    insertSQL := `INSERT INTO members (phone_number, name, visits) VALUES (?, ?, ?)`
    _, err = db.Exec(insertSQL, m.PhoneNumber, m.Name, m.Visits)
    if err != nil {
        // This usually happens if the phone number already exists
        http.Error(w, "Failed to add member: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Send a success response
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Member '%s' added successfully!", m.Name)
}

// Function to retrieve a member by phone number
func getMember(w http.ResponseWriter, r *http.Request) {
    // Only allow GET requests
    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the phone number from the URL query parameters
    phoneNumber := r.URL.Query().Get("phone_number")
    if phoneNumber == "" {
        http.Error(w, "Missing 'phone_number' parameter", http.StatusBadRequest)
        return
    }

    var m Member
    // Query the database for the member
    // We use QueryRow because we expect exactly one result (or none)
    query := `SELECT phone_number, name, visits FROM members WHERE phone_number = ?`
    row := db.QueryRow(query, phoneNumber)

    // Scan the result into our struct variables
    // We must pass pointers (&) so Scan can fill them with data
    err := row.Scan(&m.PhoneNumber, &m.Name, &m.Visits)
    if err == sql.ErrNoRows {
        http.Error(w, "Member not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Send the member data back as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(m)
}

// Function to update an existing member
func updateMember(w http.ResponseWriter, r *http.Request) {
    // Allow PUT (standard for updates) or POST
    if r.Method != http.MethodPut {
        http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var m Member
    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }

    // SQL Update statement
    updateSQL := `UPDATE members SET name = ?, visits = ? WHERE phone_number = ?`
    result, err := db.Exec(updateSQL, m.Name, m.Visits, m.PhoneNumber)
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if any row was actually updated
    rowsAffected, _ := result.RowsAffected()

    // If 0 rows were affected, it means the phone number doesn't exist in the DB
    if rowsAffected == 0 {
        http.Error(w, "Member not found. Please create the member instead.", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Member %d updated successfully!", m.PhoneNumber)
}

// Function to delete a member
func deleteMember(w http.ResponseWriter, r *http.Request) {
    // Only allow DELETE method
    if r.Method != http.MethodDelete {
        http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the phone number from the URL query parameters
    phoneNumber := r.URL.Query().Get("phone_number")
    if phoneNumber == "" {
        http.Error(w, "Missing 'phone_number' parameter", http.StatusBadRequest)
        return
    }

    // Execute the DELETE statement
    // If the number exists, it is deleted.
    // If the number does NOT exist, nothing happens, and no error is returned.
    deleteSQL := `DELETE FROM members WHERE phone_number = ?`
    _, err := db.Exec(deleteSQL, phoneNumber)
    if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Send success response
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Member deleted successfully (if they existed)")
}