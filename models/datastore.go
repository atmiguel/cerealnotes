package models

import (
	"database/sql"
	"time"
)

// ConnectToDatabase also pings the database to ensure a working connection.
func ConnectToDatabase(databaseUrl string, retry int) (*DB, error) {

	var tempDb *sql.DB
	var err error

	for attempt := 0; attempt <= retry; attempt++ {
		tempDb, err = sql.Open("postgres", databaseUrl)
		if err == nil {
			break
		}
		if err = tempDb.Ping(); err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return &DB{tempDb}, nil
}

type Datastore interface {
	// User Actions
	AuthenticateUserCredentials(*EmailAddress, string) error
	GetIdForUserWithEmailAddress(*EmailAddress) (UserId, error)
	StoreNewUser(string, *EmailAddress, string) error
	GetAllUsersById() (UserMap, error)

	// Cateogry Actions
	StoreNewNoteCategoryRelationship(NoteId, Category) error
	UpdateNoteCategory(NoteId, Category) error
	DeleteNoteCategory(NoteId) error
	GetNoteCategory(NoteId) (Category, error)

	// Note Actions
	GetUsersNotes(UserId) (NoteMap, error)
	DeleteNoteById(NoteId) error
	GetMyUnpublishedNotes(UserId) (NoteMap, error)
	StoreNewNote(*Note) (NoteId, error)
	GetAllPublishedNotesVisibleBy(UserId) (map[int64]NoteMap, error)
	GetNoteById(NoteId) (*Note, error)
	UpdateNoteContent(NoteId, string) error

	// Publication Actions
	PublishNotes(UserId) error
	StoreNewPublication(*Publication) (PublicationId, error)
}

type DB struct {
	*sql.DB
}
