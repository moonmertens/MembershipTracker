package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"net/http"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
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
func SendWhatsApp(phone int, text string) {
	if waClient == nil || !waClient.IsConnected() {
		fmt.Println("WhatsApp not connected")
		return
	}

	jid := types.NewJID(strconv.Itoa(phone), types.DefaultUserServer)

	_, err := waClient.SendMessage(context.Background(), jid, &waE2E.Message{
		Conversation: &text,
	})
	if err != nil {
		fmt.Println("Error sending WhatsApp:", err)
	}
}



