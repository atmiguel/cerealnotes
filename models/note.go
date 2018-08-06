package models

import (
	"encoding/json"
	"time"
)

type NoteId int64

type NoteType int

const (
	UNCATEGORIZED NoteType = iota
	MARGINALIA
	META
	QUESTIONS
	PREDICTIONS
)

var noteTypeStrings = [...]string{
	"uncategorized",
	"marginalia",
	"meta",
	"questions",
	"predictions",
}

func (noteType NoteType) String() string {
	if noteType < UNCATEGORIZED || noteType > PREDICTIONS {
		return "Unknown"
	}

	return noteTypeStrings[noteType]
}

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Type         NoteType  `json:"type"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
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
