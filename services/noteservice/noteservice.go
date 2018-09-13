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

func GetAllPublishedNotes() ([]*models.Note, error) {

	noteData, err := databaseutil.GetAllPublishedNotes()
	if err != nil {
		return nil, err
	}

	var notes []*models.Note = make([]*models.Note, len(noteData), len(noteData))

	for index, noteDatum := range noteData {
		notes[index] = noteDataToNote(noteDatum)
	}

	return notes, nil
}

func noteDataToNote(noteDatum *databaseutil.NoteData) *models.Note {
	return &models.Note{
		Id:           noteDatum.Id,
		AuthorId:     models.UserId(noteDatum.AuthorId),
		Content:      noteDatum.Content,
		CreationTime: noteDatum.CreationTime,
	}
}

func GetMyUnpublishedNotes(userId models.UserId) ([]*models.Note, error) {

	noteData, err := databaseutil.GetMyUnpublishedNotes(int64(userId))
	if err != nil {
		return nil, err
	}

	var notes []*models.Note = make([]*models.Note, len(noteData), len(noteData))

	for index, noteDatum := range noteData {
		notes[index] = noteDataToNote(noteDatum)
	}

	return notes, nil
}

func GetNoteById(id int64) (*models.Note, error) {
	noteData, err := databaseutil.GetNote(id)

	if err != nil {
		return nil, err
	}

	return noteDataToNote(noteData), nil
}

func DeleteNoteById(id int64) error {
	return databaseutil.DeleteNote(id)
}

func UpdateContent(id int64, content string) error {
	// TODO Do we allow content updates to published notes?
	return databaseutil.UpdateNoteContent(id, content)
}

func UpdateCategory(id int64, category models.Category) error {
	// TODO Do we allow category updates to published notes?
	return databaseutil.UpdateNoteCategory(id, category.String())
}

func DeleteCategory(id int64) error {
	if _, err := GetCategory(id); err != nil {
		if err == databaseutil.QueryResultContainedNoRowsError {
			return nil
		}

		return err
	}
	return databaseutil.DeleteNoteCategory(id)
}

func GetCategory(id int64) (models.Category, error) {
	categoryString, err := databaseutil.GetNoteCategory(id)
	if err != nil {
		return models.MARGINALIA, err
	}

	return models.DeserializeCategory(categoryString)
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
