/*
Package noteservice handles interactions with database layer.
*/
package publicationservice

import (
	"errors"

	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
)

var PublicationIdNotSet error = errors.New("The PublicationId was not set")
var PublicationAuthorDiffersFromUserAuthor error = errors.New("A note under this publication has a different author than the publication")

func CreateAndPublishNotes(userId models.UserId, notes []*models.Note) error {
	publication := models.CreateNewPublication(userId)
	if err := StoreNewPublication(publication); err != nil {
		return err
	}

	if err := PublishNotes(publication, notes); err != nil {
		return err
	}

	return nil
}

func StoreNewPublication(publication *models.Publication) error {
	id, err := databaseutil.StoreNewPublication(int64(publication.AuthorId), publication.CreationTime)
	if err != nil {
		return err
	}

	publication.Id = id

	if publication.Id < 0 {
		return PublicationIdNotSet
	}

	return nil
}

func PublishNotes(publication *models.Publication, notes []*models.Note) error {
	if publication.Id < 0 {
		return PublicationIdNotSet
	}

	noteIds := make([]int64, 0, 10)

	for _, note := range notes {

		if publication.AuthorId != note.AuthorId {
			return PublicationAuthorDiffersFromUserAuthor
		}

		noteIds = append(noteIds, note.Id)
	}

	if err := databaseutil.StorePublicationNoteRelationship(publication.Id, noteIds); err != nil {
		return err
	}

	return nil
}
