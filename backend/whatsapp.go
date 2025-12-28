package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	waClient  *whatsmeow.Client
	currentQR string
)

func InitWhatsApp() {
	dbLog := waLog.Stdout("Database", "ERROR", true)

	container, err := sqlstore.New(context.Background(), "postgres", os.Getenv("DATABASE_URL"), dbLog)
	if err != nil {
		log.Fatal(err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	waClient = whatsmeow.NewClient(deviceStore, clientLog)
	waClient.AddEventHandler(eventHandler)

	if waClient.Store.ID == nil {
		// No login found
		qrChan, _ := waClient.GetQRChannel(context.Background())
		err = waClient.Connect()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					currentQR = evt.Code
					fmt.Println("New QR Code generated")
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		}()
	} else {
		err = waClient.Connect()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("WhatsApp connected!")
	}
}

func eventHandler(evt interface{}) {
	switch evt.(type) {
	case *events.Connected:
		currentQR = ""
		fmt.Println("WhatsApp connected!")
	}
}

// Function for frontend to access QR code
func getQRCode(w http.ResponseWriter, r *http.Request) {
	status := "disconnected"

	if currentQR == "" && waClient != nil && waClient.IsConnected() {
		status = "connected"
	}

	// Send the status and the QR string (if any)
	respondJSON(w, http.StatusOK, map[string]string{
		"status": status,
		"qr":     currentQR,
	})
}

// Helper to send messages
func SendWhatsApp(phone int, text string, imgData []byte) {
	if waClient == nil || !waClient.IsConnected() {
		fmt.Println("WhatsApp not connected")
		return
	}

	jid := types.NewJID(strconv.Itoa(phone), types.DefaultUserServer)

	var msg *waE2E.Message

	if len(imgData) > 0 {
		resp, err := waClient.Upload(context.Background(), imgData, whatsmeow.MediaImage)
		if err != nil {
			fmt.Printf("Error uploading image: %v\n", err)
			return
		}

		var caption *string
		if text != "" {
			caption = &text
		}

		msg = &waE2E.Message{
			ImageMessage: &waE2E.ImageMessage{
				Caption:       caption,
				URL:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				Mimetype:      proto.String("image/jpeg"),
				FileEncSHA256: resp.FileEncSHA256,
				FileSHA256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			},
		}
	} else {
		msg = &waE2E.Message{
			Conversation: &text,
		}
	}

	_, err := waClient.SendMessage(context.Background(), jid, msg)
	if err != nil {
		fmt.Println("Error sending WhatsApp:", err)
	}
}

func broadcastMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Error parsing form data")
		return
	}

	message := r.FormValue("message")
	file, _, err := r.FormFile("image")

	// Check if image retrieval failed (other than missing file)
	if err != nil && err != http.ErrMissingFile {
		respondError(w, http.StatusBadRequest, "Error retrieving image file")
		return
	}

	var fileBytes []byte
	hasImage := (file != nil)
	if file != nil {
		defer file.Close()
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Error reading image file")
			return
		}
	}

	// Validation: Must have at least text or image
	if !hasImage && message == "" {
		respondError(w, http.StatusBadRequest, "Broadcast must contain either text or an image")
		return
	}

	// Get all members
	rows, err := db.Query("SELECT phone_number FROM members")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var phone int
		if err := rows.Scan(&phone); err != nil {
			continue
		}

		// Logic to send will go here.
		count++

		// Prepend 65 to 8-digit number
		fullPhone := phone + 6500000000
		go SendWhatsApp(fullPhone, message, fileBytes)
	}

	statusMsg := fmt.Sprintf("Broadcast queued for %d members. Image attached: %v", count, hasImage)
	respondJSON(w, http.StatusOK, map[string]string{"message": statusMsg})
}
