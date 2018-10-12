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
	return 0, CannotDeserializeCategoryStringError
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

func (db *DB) GetNoteCategory(noteId NoteId) (Category, error) {

	sqlQuery := `
		SELECT category FROM note_to_category_relationship
		WHERE note_id = $1`

	var categoryString string
	if err := db.execOneResult(sqlQuery, &categoryString, int64(noteId)); err != nil {
		return 0, err
	}

	category, err := DeserializeCategory(categoryString)
	if err != nil {
		return 0, err
	}

	return category, nil
}

func (db *DB) UpdateNoteCategory(noteId NoteId, category Category) error {
	sqlQuery := `
		INSERT INTO note_to_category_relationship (note_id, category)
		VALUES ($1, $2)
		ON CONFLICT (note_id) DO 
		UPDATE SET category = ($2)
		WHERE note_to_category_relationship.note_id = ($1)`

	rowsAffected, err := db.execNoResults(sqlQuery, int64(noteId), category.String())
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return NoNoteFoundError
	}

	if rowsAffected > 1 {
		return TooManyRowsAffectedError
	}

	return nil
}

func (db *DB) DeleteNoteCategory(noteId NoteId) error {
	sqlQuery := `
		DELETE FROM note_to_category_relationship
		WHERE note_id = $1`

	num, err := db.execNoResults(sqlQuery, int64(noteId))
	if err != nil {
		return err
	}

	if num == 0 {
		return NoNoteFoundError
	}

	if num != 1 {
		return TooManyRowsAffectedError
	}

	return nil
}
