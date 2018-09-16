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

	if category < MARGINALIA || category > PREDICTIONS {
		return "Unknown"
	}

	return categoryStrings[category]
}

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}

//  DB methods

func (db *DB) StoreNewNote(
	note *Note,
) (NoteId, error) {

	authorId := int64(note.AuthorId)
	content := note.Content
	creationTime := note.CreationTime

	sqlQuery := `
		INSERT INTO note (author_id, content, creation_time)
		VALUES ($1, $2, $3)
		RETURNING id`

	rows, err := db.Query(sqlQuery, authorId, content, creationTime)
	if err != nil {
		return 0, convertPostgresError(err)
	}
	defer rows.Close()

	var noteId int64 = 0
	for rows.Next() {

		if noteId != 0 {
			return 0, QueryResultContainedMultipleRowsError
		}

		if err := rows.Scan(&noteId); err != nil {
			return 0, convertPostgresError(err)
		}
	}

	if noteId == 0 {
		return 0, QueryResultContainedNoRowsError
	}

	if err := rows.Err(); err != nil {
		return 0, convertPostgresError(err)
	}

	return NoteId(noteId), nil
}
