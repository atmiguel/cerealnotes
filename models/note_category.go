package models

import (
	"errors"
)

type NoteCategory int

const (
	MARGINALIA NoteCategory = iota
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

var CannotDeserializeNoteCategoryStringError = errors.New("String does not correspond to a Note Category")

func DeserializeNoteCategory(input string) (NoteCategory, error) {
	for i := 0; i < len(categoryStrings); i++ {
		if input == categoryStrings[i] {
			return NoteCategory(i), nil
		}
	}
	return 0, CannotDeserializeNoteCategoryStringError
}

func (category NoteCategory) String() string {

	if category < MARGINALIA || category > PREDICTIONS {
		return "Unknown"
	}

	return categoryStrings[category]
}

func (db *DB) StoreNewNoteCategoryRelationship(
	noteId NoteId,
	category NoteCategory,
) error {
	sqlQuery := `
		INSERT INTO note_to_category_relationship (note_id, category)
		VALUES ($1, $2)`

	rows, err := db.Query(sqlQuery, int64(noteId), category.String())
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}
