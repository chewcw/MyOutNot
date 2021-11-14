package main

import (
	"encoding/json"
	"fmt"
	"integration/model"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type GraphService struct {
	authz    *AuthzService
	tenantID string
}

func NewGraphService(authz *AuthzService) *GraphService {
	return &GraphService{
		authz:    authz,
		tenantID: tenantID,
	}
}

func (g *GraphService) FetchEvents() {
	if g.authz.accessToken == "" {
		log.Println("No access token, not going to fetch events for now")
		return
	}

	log.Println("Fetching events")

	if g.authz.userOid == "" {
		log.Println("User oid is empty")
	}
	req, _ := http.NewRequest(
		"GET",
		"https://graph.microsoft.com/v1.0/"+tenantID+"/users/"+g.authz.userOid+"/calendar/events",
		nil)
	req.Header.Add("Authorization", "Bearer "+g.authz.accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var events model.RawEvent
	if err := json.Unmarshal(body, &events); err != nil {
		log.Println(err)
		return
	}

	// check which event need to be added
	now := time.Now().UTC()
	lastChecked := now.Add(-checkEventsDuration)

	for _, event := range events.Value {
		createdDateTime, _ := time.Parse("2006-01-02T15:04:05.0000000Z", event.CreatedDateTime)
		if createdDateTime.Before(now) && createdDateTime.After(lastChecked) {
			// attendees
			attendees := []string{}
			for _, attendee := range event.Attendees {
				name := fmt.Sprintf("@%s", attendee.EmailAddress.Name)
				attendees = append(attendees, name)
			}

			// meeting start datetime
			startDateTime, _ := time.Parse("2006-01-02T15:04:05.0000000", event.Start.DateTime)
			startDateTimeStr := startDateTime.Format("2006-01-02T15:04:05+08:00")
			if strings.ToLower(event.Start.TimeZone) == "utc" {
				startDateTimeLocation, err := time.LoadLocation("Asia/Kuala_Lumpur")
				if err != nil {
					log.Fatal(err)
				}
				startDateTimeStr = startDateTime.In(startDateTimeLocation).Format("2006-01-02T15:04:05+08:00")
			}

			notionService := NewNotionService(
				event.Subject,
				startDateTimeStr,
				strings.Join(attendees, ","),
				"External",
			)

			notionService.CreatePage()
		}
	}
}
