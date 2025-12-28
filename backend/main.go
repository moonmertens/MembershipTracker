package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
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
	initDB()
	defer db.Close()

	go InitWhatsApp()

	mux := http.NewServeMux()
	mux.HandleFunc("/add-member", addMember)
	mux.HandleFunc("/get-member", getMember)
	mux.HandleFunc("/update-member", updateMember)
	mux.HandleFunc("/delete-member", deleteMember)
	mux.HandleFunc("/get-all-members", getAllMembers)
	mux.HandleFunc("/broadcast-message", broadcastMessage)
	mux.HandleFunc("/get-whatsapp-qr", getQRCode)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins (for now)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on port " + port)
	http.ListenAndServe(":"+port, handler)
}

func initDB() {

	var err error

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS members (
        phone_number INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        visits INTEGER DEFAULT 0
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
