package models

import (
	"encoding/json"
	"time"
)

type NoteId int64

type NoteType int

const (
	MARGINALIA NoteType = iota
	META
	QUESTIONS
	PREDICTIONS
)

const noteTypeStrings = [...]string{
	"Marginalia",
	"Meta",
	"Questions",
	"Predictions",
}

func (noteType NoteType) String() string {
	if noteType < MARGINALIA || noteType > PREDICTIONS {
		return "Unknown"
	}

	return noteTypeStrings[noteType]
}

type Note struct {
	AuthorId      UserId        `json:"authorId"`
	Type          NoteType      `json:"type"`
	Content       string        `json:"content"`
	PublicationId PublicationId `json:"publicationId"`
	CreationTime  time.Time     `json:"creationTime"`
}

func (note *Note) MarshalJSON() ([]byte, error) {
	type Alias Note

	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  note.Type.String(),
		Alias: (*Alias)(note),
	})
}
