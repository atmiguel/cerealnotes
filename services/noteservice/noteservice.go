/*
Package noteservice handles interactions with database layer.
*/
package noteservice

import (
	"errors"

	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
)

var NoteIdIsNotValid error = errors.New("Email address already in use")

func StoreNewNote(
	note *models.Note,
) error {
	databaseutil.StoreNewNote(int64(note.AuthorId), note.Content, note.CreationTime)
	return nil
}

func StoreNoteCategoryRelationship(
	note *models.Note,
	category models.Category,
) error {
	return nil
}
