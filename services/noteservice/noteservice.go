/*
Package noteservice handles interactions with database layer.
*/
package noteservice

import (
	"errors"

	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
)

var NoteIdNotSet error = errors.New("The NoteId was not set")

func StoreNewNote(
	note *models.Note,
) error {

	id, err := databaseutil.StoreNewNote(int64(note.AuthorId), note.Content, note.CreationTime)
	if err != nil {
		return err
	}

	note.Id = id

	if note.Id < 0 {
		return NoteIdNotSet
	}

	return nil
}

func StoreNoteCategoryRelationship(
	note *models.Note,
	category models.Category,
) error {

	if note.Id < 0 {
		return NoteIdNotSet
	}

	if err := databaseutil.StoreNoteCategoryRelationship(int64(note.Id), category.String()); err != nil {
		return err
	}

	return nil
}
