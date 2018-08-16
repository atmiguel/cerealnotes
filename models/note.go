package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type NoteId int64

type NoteType int

var InvalidNoteTypeError = errors.New("Note type does not exist")

var UnDeserializeableNoteTypeStringError = errors.New("String does not correspond to a NoteType")

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

func DeserializeNoteType(input string) (NoteType, error) {
	for i := 0; i < len(noteTypeStrings); i++ {
		if strings.ToLower(input) == strings.ToLower(noteTypeStrings[i]) {
			return NoteType(i), nil
		}
	}
	return UNCATEGORIZED, UnDeserializeableNoteTypeStringError
}

func (noteType NoteType) String() (string, error) {
	if noteType < UNCATEGORIZED || noteType > PREDICTIONS {
		return "", InvalidNoteTypeError
	}

	return noteTypeStrings[noteType], nil
}

func CreateNewNote(userId UserId, content string, noteType NoteType) *Note {
	return &Note{
		AuthorId:     userId,
		Content:      content,
		Type:         noteType,
		CreationTime: time.Now().UTC(),
	}
}

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Type         NoteType  `json:"type"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}

func (note *Note) MarshalJSON() ([]byte, error) {
	type Alias Note

	if notetype, err := note.Type.String(); err == nil {
		return json.Marshal(&struct {
			Type string `json:"type"`
			*Alias
		}{
			Type:  notetype,
			Alias: (*Alias)(note),
		})
	}

	return nil, InvalidNoteTypeError

}
