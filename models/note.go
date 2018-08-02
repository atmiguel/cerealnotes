package models

import (
	"encoding/json"
	"strings"
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

var noteTypeStrings = [...]string{
	"Marginalia",
	"Meta",
	"Questions",
	"Predictions",
}

func DeserializeNoteType(input string) NoteType {
	for i := 0; i < len(noteTypeStrings); i++ {
		if strings.ToLower(input) == strings.ToLower(noteTypeStrings[i]) {
			return NoteType(i)
		}
	}
	return NoteType(-1)
}

func (noteType NoteType) String() string {
	if noteType < MARGINALIA || noteType > PREDICTIONS {
		return ""
	}

	return noteTypeStrings[noteType]
}

func CreateNewNote(userId UserId, content string, noteType NoteType) *Note {
	note := new(Note)
	note.AuthorId = userId
	note.Content = content
	note.Type = noteType
	note.CreationTime = time.Now()
	note.PublicationId = -1

	return note
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
