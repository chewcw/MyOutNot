package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to Lalaland")
	})

	// create
	engine.POST("/create", func(c *gin.Context) {
		h, err := json.Marshal(c.Request.Header)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		var header map[string][]string
		if err := json.Unmarshal(h, &header); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if len(header["Api_key"]) == 0 {
			c.JSON(http.StatusUnauthorized, "Please provide Api_key header")
			return
		}

		if header["Api_key"][0] == apiKey {
			var res map[string]interface{}
			body, _ := ioutil.ReadAll(c.Request.Body)
			if len(body) > 0 {
				err = json.Unmarshal(body, &res)
				if err != nil {
					c.JSON(http.StatusBadRequest, err.Error())
					return
				}

				// meeting start datetime
				startDateTime, _ := time.Parse("2006-01-02T15:04:05+00:00", res["LocalStartDateTime"].(string))
				startDateTimeLocation, err := time.LoadLocation("Asia/Kuala_Lumpur")
				if err != nil {
					log.Fatal(err)
				}
				startDateTimeStr := startDateTime.In(startDateTimeLocation).Format("2006-01-02T15:04:05+08:00")

				// attendees
				splitAttendees := strings.Split(res["Attendees"].(string), ",")
				attendees := []string{}
				for i := 0; i < len(splitAttendees); i++ {
					if i > 4 {
						continue
					}
					attendees = append(attendees, splitAttendees[i])
				}

				notionService := NewNotionService(
					res["Title"].(string),
					startDateTimeStr,
					strings.Join(attendees, ","),
					externalOrInternalV2(res["Organizer"].(string)),
				)

				notionService.CreatePage()

				c.JSON(http.StatusOK, "")
				return
			}
			c.JSON(http.StatusNoContent, "")
		} else {
			c.JSON(http.StatusUnauthorized, "Api_key doesn't match")
			return
		}
	})

	return engine
}

func externalOrInternalV2(organizer string) string {
	if strings.Contains(
		organizer,
		externalOrganizerEmail) {
		return "External"
	}
	return "Internal"
}
