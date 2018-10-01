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

	rows, err := db.Query(sqlQuery, authorId, content, creationTime)
	if err != nil {
		return 0, convertPostgresError(err)
	}
	defer rows.Close()

	var noteId int64 = 0
	for rows.Next() {

		if noteId != 0 {
			return 0, QueryResultContainedMultipleRowsError
		}

		if err := rows.Scan(&noteId); err != nil {
			return 0, convertPostgresError(err)
		}
	}

	if noteId == 0 {
		return 0, QueryResultContainedNoRowsError
	}

	if err := rows.Err(); err != nil {
		return 0, convertPostgresError(err)
	}

	return NoteId(noteId), nil
}
