package model

type NotionPage struct {
	Parent     Parent     `json:"parent"`
	Properties []Property `json:"properties"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type Property struct {
	Title           Title           `json:"title"`
	MeetingDateTime MeetingDateTime `json:"Meeting date & time"`
	Attendees       Attendees       `json:"attendees"`
	Type            Type            `json:"type"`
}

type Title struct {
	Title []Text `json:"title"`
}

type Text struct {
	Text Content `json:"text"`
}

type Content struct {
	Content string `json:"content"`
}

type MeetingDateTime struct {
	Date Date `json:"date"`
}

type Date struct {
	Start string `json:"start"`
}

type Attendees struct {
	RichTexts []Text `json:"rich_text"`
}

type Type struct {
	Select Select `json:"select"`
}

type Select struct {
	Name string `json:"name"`
}
