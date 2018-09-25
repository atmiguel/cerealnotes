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
var NoteAlreadyContainsCategoryError = errors.New("NoteId already has a category stored for it")

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

	if _, err := db.execNoResults(sqlQuery, int64(noteId), category.String()); err != nil {
		if err == UniqueConstraintError {
			return NoteAlreadyContainsCategoryError
		}
		return err
	}

	return nil
}
