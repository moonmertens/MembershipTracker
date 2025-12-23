package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/rs/cors"
    _ "modernc.org/sqlite"
)

// Member struct is shared across files in package main
type Member struct {
    PhoneNumber int    `json:"phone_number"`
    Name        string `json:"name"`
    Visits      int    `json:"visits"`
}

// db variable is shared across files in package main
var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite", "./members.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    initDB()

    mux := http.NewServeMux()
    mux.HandleFunc("/add-member", addMember)
    mux.HandleFunc("/get-member", getMember)
    mux.HandleFunc("/update-member", updateMember)
    mux.HandleFunc("/delete-member", deleteMember)

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Allow all origins (for now)
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })

    handler := c.Handler(mux)

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}

func initDB() {
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS members (
        phone_number INTEGER PRIMARY KEY, 
        name TEXT,
        visits INTEGER
    );
    `
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatalf("Database creation failed: %q: %s\n", err, sqlStmt)
    }
}