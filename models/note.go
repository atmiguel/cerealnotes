package models

import (
	"errors"
	"time"
)

type NoteId int64

type Category int

const (
	MARGINALIA Category = iota
	META
	QUESTIONS
	PREDICTIONS
)

var categoryStrings = [...]string{
	"marginalia",
	"meta",
	"questions",
	"predictions",
}

var CannotDeserializeCategoryStringError = errors.New("String does not correspond to a Note Category")

func DeserializeCategory(input string) (Category, error) {
	for i := 0; i < len(categoryStrings); i++ {
		if input == categoryStrings[i] {
			return Category(i), nil
		}
	}
	return MARGINALIA, CannotDeserializeCategoryStringError
}

func (category Category) String() string {
	return categoryStrings[category]
}

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}
