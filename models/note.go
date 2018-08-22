package models

import (
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
	"marginalia",
	"meta",
	"questions",
	"predictions",
}

func (noteType NoteType) String() string {
	if noteType < MARGINALIA || noteType > PREDICTIONS {
		return "Unknown"
	}

	return noteTypeStrings[noteType]
}

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}
