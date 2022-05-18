package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	notionNewPageUrl = "https://api.notion.com/v1/pages"
)

type NotionService struct {
	Title              string
	LocalStartDateTime string
	Attendees          string
	Type               string
	DatabaseID         string
	Secret             string
}

func NewNotionService(title, localStartDateTime, attendees, typ string) *NotionService {
	return &NotionService{
		Title:              title,
		LocalStartDateTime: localStartDateTime,
		Attendees:          attendees,
		Type:               typ,
		Secret:             secret,
		DatabaseID:         databaseID,
	}
}

func (n *NotionService) CreatePage() {
	log.Printf("Creating notion page titled \"%s\"", n.Title)

	postBody := fmt.Sprintf(`
{
	"parent": {
		"database_id": "%s"
	},
	"properties": {
		"title": {
			"title": [
				{
					"text": {
						"content": "%s"
					}
				}
			]
		},
		"Meeting date & time": {
			"date": {
				"start": "%s"
			}
		},
		"Attendees": {
			"rich_text": [
				{
					"text": {
						"content": "%s"
					}
				}
			]
		},
		"Type": {
			"select": {
				"name": "%s"
			}
		}
	}
}
	`, n.DatabaseID, n.Title, n.LocalStartDateTime, n.Attendees, n.Type)

	var jsonStr = []byte(postBody)
	req, _ := http.NewRequest("POST", notionNewPageUrl, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+n.Secret)
	req.Header.Add("Notion-Version", "2021-08-16")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	res, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf("Create notion page returned %d, %s", resp.StatusCode, res)
		return
	}

	log.Printf("Created notion page titled %s", n.Title)
}
