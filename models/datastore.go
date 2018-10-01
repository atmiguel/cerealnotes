package models

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// UniqueConstraintError is returned when a uniqueness constraint is violated during an insert.
var UniqueConstraintError = errors.New("postgres: unique constraint violation")

// QueryResultContainedMultipleRowsError is returned when a query unexpectedly returns more than one row.
var QueryResultContainedMultipleRowsError = errors.New("query result unexpectedly contained multiple rows")

// QueryResultContainedNoRowsError is returned when a query unexpectedly returns no rows.
var QueryResultContainedNoRowsError = errors.New("query result unexpectedly contained no rows")

// ConnectToDatabase also pings the database to ensure a working connection.
func ConnectToDatabase(databaseUrl string) (*DB, error) {
	tempDb, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	if err := tempDb.Ping(); err != nil {
		return nil, err
	}

	return &DB{tempDb}, nil
}

type Datastore interface {
	StoreNewNote(*Note) (NoteId, error)
	StoreNewNoteCategoryRelationship(NoteId, NoteCategory) error
	StoreNewUser(string, *EmailAddress, string) error
	AuthenticateUserCredentials(*EmailAddress, string) error
	GetIdForUserWithEmailAddress(*EmailAddress) (UserId, error)
}

type DB struct {
	*sql.DB
}

func convertPostgresError(err error) error {
	const uniqueConstraintErrorCode = "23505"

	if postgresErr, ok := err.(*pq.Error); ok {
		if postgresErr.Code == uniqueConstraintErrorCode {
			return UniqueConstraintError
		}
	}

	return err
}
