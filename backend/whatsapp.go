package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
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
