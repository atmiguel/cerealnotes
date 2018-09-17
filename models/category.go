package models

import (
	"errors"
)

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

func (db *DB) StoreNewNoteCategoryRelationship(
	noteId NoteId,
	category Category,
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
