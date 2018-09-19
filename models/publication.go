package models

import (
	"errors"
	"strconv"
	"time"
)

type PublicationId int64

type Publication struct {
	AuthorId     UserId    `json:"authorId"`
	CreationTime time.Time `json:"creationTime"`
}

var NoNotesToPublishError = errors.New("There are no unpublished notes to publish")

func (db *DB) PublishNotes(userId UserId) error {

	myUnpublishedNotes, err := db.GetMyUnpublishedNotes(userId)
	if err != nil {
		return err
	}

	if len(myUnpublishedNotes) == 0 {
		return NoNotesToPublishError
	}

	publicationId, err := db.StoreNewPublication(&Publication{AuthorId: userId, CreationTime: time.Now().UTC()})

	sqlQuery := `
		INSERT INTO note_to_publication_relationship (publication_id, note_id) 
		VALUES ($1, $2)`

	noteIds := make([]int64, len(myUnpublishedNotes))

	{
		i := 0
		for noteId := range myUnpublishedNotes {
			noteIds[i] = int64(noteId)
			i++
		}
	}

	values := make([]interface{}, len(noteIds)*2, len(noteIds)*2)
	{
		values[0] = publicationId
		values[1] = noteIds[0]
		for index, noteId := range noteIds {
			if index == 0 {
				continue
			}
			sqlQuery += ", ($" + strconv.Itoa(2*index+1) + ", $" + strconv.Itoa((2*index)+2) + ")"
			values[2*index] = publicationId
			values[(2*index)+1] = noteId
		}
	}

	rowAffected, err := db.execNoResults(sqlQuery, values...)
	if err != nil {
		return err
	}

	if rowAffected < 1 {
		return errors.New("No values were inserted")
	}

	return nil

}

func (db *DB) StoreNewPublication(publication *Publication) (PublicationId, error) {

	sqlQuery := `
		INSERT INTO publication (author_id, creation_time)
		VALUES ($1, $2)
		RETURNING id`

	var publicationId int64 = 0
	if err := db.execOneResult(sqlQuery, &publicationId, int64(publication.AuthorId), publication.CreationTime); err != nil {
		return 0, err
	}

	return PublicationId(publicationId), nil
}
