package models

import "time"

type NoteId int64

type NoteType int

const (
	MARGINALIA NoteType = iota
	META
	QUESTIONS
	PREDICTIONS
)

type Note struct {
	AuthorId      UserId        `json:"authorId"`
	Type          NoteType      `json:"type"`
	Content       string        `json:"content"`
	PublicationId PublicationId `json:"publicationId"`
	CreationTime  time.Time     `json:"creationTime"`
}
