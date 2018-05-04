/*
Package databaseutil abstracts away details about sql and postgres.

These functions only accept and return primitive types.
*/
package databaseutil

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var db *sql.DB

// UniqueConstraintError is returned when a uniqueness constraint is violated during an insert.
var UniqueConstraintError = errors.New("postgres: unique constraint violation")

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

func InsertIntoUsersTable(
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
	var row *sql.Row
	{
		sqlQuery := `
			SELECT password FROM users
			WHERE email_address = $1`

		row = db.QueryRow(sqlQuery, emailAddress)
	}

	var password []byte
	if err := row.Scan(&password); err != nil {
		return nil, err
	}

	return password, nil
}

func GetIdForUserWithEmailAddress(emailAddress string) (int64, error) {
	var row *sql.Row
	{
		sqlQuery := `
			SELECT id FROM users
			WHERE email_address = $1`

		row = db.QueryRow(sqlQuery, emailAddress)
	}

	var userId int64
	if err := row.Scan(&userId); err != nil {
		return 0, err
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
