package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var secret = os.Getenv("NOTION_SECRET")
var databaseID = os.Getenv("NOTION_DATABASE_ID")
var externalOrganizerEmail = os.Getenv("EXTERNAL_ORGANIZER_EMAIL")
var apiKey = os.Getenv("API_KEY")
var checkEventsDuration time.Duration

func init() {
	if apiKey == "" {
		log.Fatal("apiKey env is empty")
	}

	if secret == "" || databaseID == "" {
		log.Fatal("Notion secret or databaseID env is empty")
	}

	// checkEventsDuration env fallback
	duration := 5
	if value, ok := os.LookupEnv("CHECK_EVENTS_DURATION"); ok {
		var err error
		duration, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal("Wrong setting on checkEventsDuration env")
		}
	}
	checkEventsDuration = time.Duration(duration) * time.Minute

	msg := fmt.Sprintf(`Settings:
secret: %s,
databaseID: %s,
apiKey: %s,
checkEventDuration: %s`,
		secret,
		databaseID,
		apiKey,
		checkEventsDuration,
	)
	log.Println(msg)
}
