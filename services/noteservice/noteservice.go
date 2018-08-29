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

func GetNoteById(id int64) (*models.Note, error) {
	noteData, err := databaseutil.GetNote(id)

	if err != nil {
		return nil, err
	}

	return &models.Note{
			Id:           noteData.Id,
			AuthorId:     models.UserId(noteData.AuthorId),
			Content:      noteData.Content,
			CreationTime: noteData.CreationTime,
		},
		nil
}

func DeleteNoteById(id int64) error {
	return databaseutil.DeleteNote(id)
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
