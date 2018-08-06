/*
Package databaseutil abstracts away details about sql and postgres.

These functions only accept and return primitive types.
*/
package databaseutil

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/lib/pq"
)

var db *sql.DB

// UniqueConstraintError is returned when a uniqueness constraint is violated during an insert.
var UniqueConstraintError = errors.New("postgres: unique constraint violation")

// QueryResultContainedMultipleRowsError is returned when a query unexpectedly returns more than one row.
var QueryResultContainedMultipleRowsError = errors.New("query result unexpectedly contained multiple rows")

// QueryResultContainedNoRowsError is returned when a query unexpectedly returns no rows.
var QueryResultContainedNoRowsError = errors.New("query result unexpectedly contained no rows")

// ConnectToDatabase also pings the database to ensure a working connection.
func ConnectToDatabase(databaseUrl string) error {
	{
		tempDb, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}

		db = tempDb
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func InsertIntoUserTable(
	displayName string,
	emailAddress string,
	password []byte,
	creationTime time.Time,
) error {
	sqlQuery := `
		INSERT INTO users (display_name, email_address, password, creation_time)
		VALUES ($1, $2, $3, $4)`

	rows, err := db.Query(sqlQuery, displayName, emailAddress, password, creationTime)
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}

func GetPasswordForUserWithEmailAddress(emailAddress string) ([]byte, error) {
	sqlQuery := `
		SELECT password FROM users
		WHERE email_address = $1`

	rows, err := db.Query(sqlQuery, emailAddress)
	if err != nil {
		return nil, convertPostgresError(err)
	}
	defer rows.Close()

	var password []byte
	for rows.Next() {
		if password != nil {
			return nil, QueryResultContainedMultipleRowsError
		}

		if err := rows.Scan(&password); err != nil {
			return nil, err
		}
	}

	if password == nil {
		return nil, QueryResultContainedNoRowsError
	}

	return password, nil
}

func StoreNewNote(note *models.Note) error {
	sqlQuery := `
		INSERT INTO notes (author_id, type, content, publication_id, creation_time)
		VALUES ($1, $2, $3, $4, $5)`

	var cleanPublicationId sql.NullInt64
	var cleanNoteType sql.NullString

	publicationId := int64(note.PublicationId)
	if publicationId < 1 {
		cleanPublicationId = sql.NullInt64{Int64: 0, Valid: false}
	} else {
		cleanPublicationId = sql.NullInt64{Int64: publicationId, Valid: true}
	}

	noteType := note.Type.String()
	if len(noteType) == 0 {
		cleanNoteType = sql.NullString{String: "", Valid: false}
	} else {
		cleanNoteType = sql.NullString{String: strings.ToLower(noteType), Valid: true}
	}

	rows, err := db.Query(sqlQuery, int64(note.AuthorId), cleanNoteType, note.Content, cleanPublicationId, note.CreationTime)
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}

func convertRowToNote(rows *sql.Rows) models.Note {
	return models.Note{}
}

func GetAllPublishedNotes() ([]models.Note, error) {
	sqlQuery := `
		SELECT * FROM notes
		WHERE publication_id != NULL
		ORDER BY creation_time`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, convertPostgresError(err)
	}

	defer rows.Close()
	notes := make([]models.Note, 0)

	for rows.Next() {
		notes = append(notes, convertRowToNote(rows))
	}

	return notes, nil
}

func GetIdForUserWithEmailAddress(emailAddress string) (int64, error) {
	sqlQuery := `
		SELECT id FROM users
		WHERE email_address = $1`

	rows, err := db.Query(sqlQuery, emailAddress)
	if err != nil {
		return 0, convertPostgresError(err)
	}
	defer rows.Close()

	var userId int64
	for rows.Next() {
		if userId != 0 {
			return 0, QueryResultContainedMultipleRowsError
		}

		if err := rows.Scan(&userId); err != nil {
			return 0, err
		}
	}

	if userId == 0 {
		return 0, QueryResultContainedNoRowsError
	}

	return userId, nil
}

// PRIVATE

func convertPostgresError(err error) error {
	const uniqueConstraintErrorCode = "23505"

	if postgresErr, ok := err.(*pq.Error); ok {
		if postgresErr.Code == uniqueConstraintErrorCode {
			return UniqueConstraintError
		}
	}

	return err
}
