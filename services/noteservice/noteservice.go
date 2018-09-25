/*
Package noteservice handles interactions with database layer.
*/
package noteservice

import (
	"encoding/json"
	"fmt"

	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
)

func StoreNewNote(
	note *models.Note,
) (models.NoteId, error) {

	id, err := databaseutil.InsertNewNote(int64(note.AuthorId), note.Content, note.CreationTime)
	if err != nil {
		return models.NoteId(0), err
	}

	return models.NoteId(id), nil
}

func StoreNewNoteCategoryRelationship(
	noteId models.NoteId,
	category models.Category,
) error {
	if err := databaseutil.InsertNoteCategoryRelationship(int64(noteId), category.String()); err != nil {
		return err
	}

	return nil
}

type NotesById map[models.NoteId]*models.Note

func (noteMap NotesById) ToJson() ([]byte, error) {
	// json doesn't support int indexed maps
	notesByIdString := make(map[string]models.Note, len(noteMap))

	for id, note := range noteMap {
		notesByIdString[fmt.Sprint(id)] = *note
	}

	return json.Marshal(notesByIdString)
}
