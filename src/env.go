package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
)

var secret = os.Getenv("NOTION_SECRET")
var databaseID = os.Getenv("NOTION_DATABASE_ID")
var tenantID = os.Getenv("AZURE_TENANT_ID")
var clientID = os.Getenv("AAD_CLIENT_ID")
var clientSecret = os.Getenv("AAD_CLIENT_SECRET")
var redirectURL = os.Getenv("AAD_REDIRECT_URL")
var azureTableConnStr = os.Getenv("AZ_TABLE_CONN_STR")
var azureTableName = os.Getenv("AZ_TABLE_NAME")
var azureTablePartitionKey = os.Getenv("AZ_TABLE_PARTITION_KEY")
var azureTableRowKey = os.Getenv("AZ_TABLE_ROW_KEY")
var app = os.Getenv("APP")
var externalOrganizerEmail = os.Getenv("EXTERNAL_ORGANIZER_EMAIL")
var azureTableFullMetadata = storage.FullMetadata
var checkEventsDuration time.Duration

func init() {
	if tenantID == "" || clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Fatal("tenantID, clientID, clientSecret or redirectURL env is empty")
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
tenantID: %s,
clientID: %s,
clientSecret: %s,
redirectURL: %s,
checkEventDuration: %s`,
		secret,
		databaseID,
		tenantID,
		clientID,
		clientSecret,
		redirectURL,
		checkEventsDuration,
	)
	log.Println(msg)
}
