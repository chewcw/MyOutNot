package main

import (
	"encoding/json"
	"fmt"
	"integration/model"
	"log"
	"strings"
	"time"
)

type GraphService struct {
}

func NewGraphService() *GraphService {
	return &GraphService{}
}

func (g *GraphService) FetchEvents(body []byte) {
	log.Println("Fetching events")

	if len(body) == 0 {
		log.Println("No event")
		return
	}

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
			i := 0
			for _, attendee := range event.Attendees {
				if i > 4 {
					continue
				}
				name := fmt.Sprintf("@%s", attendee.EmailAddress.Name)
				attendees = append(attendees, name)
				i++
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
				externalOrInternal(event.Organizer),
			)

			notionService.CreatePage()
		}
	}
}

func externalOrInternal(organizer model.Attendee) string {
	if strings.Contains(
		organizer.EmailAddress.Address,
		externalOrganizerEmail) {
		return "External"
	}
	return "Internal"
}
