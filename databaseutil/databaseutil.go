/*
Package databaseutil provides functions to be run against the database.


These functions are simple wrappers around the databse accepting and returning
primitive types.
*/
package databaseutil

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var db *sql.DB

// UniqueConstraintError is returned when a uniqueness constraint is violated
// when trying to insert into the table
var UniqueConstraintError = errors.New("postgres: unique constraint violation")

// ConnectToDatabase connects to the database and pings the database to make
// sure that the connection works.
func ConnectToDatabase(databaseUrl string) error {
	{
		tempDb, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}

		db = tempDb
	}

	// Quickly test if the connection to the database worked.
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

// InsertIntoUsersTable inserts the given user information into the database
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

	if err = rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}

// GetPasswordForUserWithEmailAddress given an email address returns the password
// as []byte
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

// GetIdForUserWithEmailAddress returns the user id assosiated with the given
// email address
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
		return -1, err
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
