package models

import (
	"errors"
	"time"
)

type NoteId int64

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}

var NoNoteFoundError = errors.New("No note with that information could be found")

//  DB methods

func (db *DB) StoreNewNote(
	note *Note,
) (NoteId, error) {

	authorId := int64(note.AuthorId)
	content := note.Content
	creationTime := note.CreationTime

	sqlQuery := `
		INSERT INTO note (author_id, content, creation_time)
		VALUES ($1, $2, $3)
		RETURNING id`

	var noteId int64 = 0
	if err := db.execOneResult(sqlQuery, &noteId, authorId, content, creationTime); err != nil {
		return 0, err
	}
	return NoteId(noteId), nil
}

func (db *DB) GetUsersNotes(userId UserId) (NoteMap, error) {
	noteMap := make(map[NoteId]*Note)

	{
		sqlQuery := `
			SELECT id, author_id, content, creation_time FROM note
			WHERE author_id = $1`
		rows, err := db.Query(sqlQuery, int64(userId))
		if err != nil {
			return nil, convertPostgresError(err)
		}
		defer rows.Close()

		for rows.Next() {
			var tempId int64
			tempNote := &Note{}
			if err := rows.Scan(&tempId, &tempNote.AuthorId, &tempNote.Content, &tempNote.CreationTime); err != nil {
				return nil, convertPostgresError(err)
			}

			noteMap[NoteId(tempId)] = tempNote
		}
	}

	return noteMap, nil
}

func (db *DB) GetAllPublishedNotesVisibleBy(userId UserId) (NoteMap, error) {
	return nil, errors.New("Not implemented")
}

func (db *DB) GetMyUnpublishedNotes(userId UserId) (NoteMap, error) {
	return nil, errors.New("Not implimented")
}

func (db *DB) DeleteNoteById(noteId NoteId) error {
	sqlQuery := `
		DELETE FROM note
		WHERE id = $1`

	num, err := db.execNoResults(sqlQuery, int64(noteId))
	if err != nil {
		return err
	}

	if num == 0 {
		return NoNoteFoundError
	}

	if num != 1 {
		return errors.New("Somewhere we more than 1 note was deleted")
	}

	return nil
}
