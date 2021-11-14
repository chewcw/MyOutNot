package model

type RawEvent struct {
	Value []Event `json:"value"`
}

type Event struct {
	IsReminderOn    bool       `json:"isReminderOn"`
	Subject         string     `json:"subject"`
	Organizer       Attendee   `json:"organizer"`
	Attendees       []Attendee `json:"attendees"`
	Start           Time       `json:"start"`
	End             Time       `json:"end"`
	CreatedDateTime string     `json:"createdDateTime"`
	WebLink         string     `json:"webLink"`
}

type Attendee struct {
	Type         string       `json:"type"`
	EmailAddress EmailAddress `json:"emailAddress"`
}

type EmailAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Time struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}
