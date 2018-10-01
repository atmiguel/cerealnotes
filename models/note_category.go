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
var NoteAlreadyContainsCategoryError = errors.New("NoteId already has a category stored for it")

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

	if _, err := db.execNoResults(sqlQuery, int64(noteId), category.String()); err != nil {
		if err == UniqueConstraintError {
			return NoteAlreadyContainsCategoryError
		}
		return err
	}

	return nil
}
