package models

import (
	"time"
)

type NoteId int64

type Note struct {
	AuthorId     UserId    `json:"authorId"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creationTime"`
}

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
