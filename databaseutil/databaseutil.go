package databaseutil

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

var db *sql.DB

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

func InsertIntoUsersTable(
	displayName string,
	emailAddress string,
	password []byte,
	creationTime time.Time,
) error {
	var row *sql.Row
	{
		sqlQuery := `
			INSERT INTO users (display_name, email_address, password, creation_time)
			VALUES ($1, $2, $3, $4)`

		row = db.QueryRow(sqlQuery, displayName, emailAddress, password, creationTime)
	}

	if err := row.Scan(); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
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

func GetUserIdFromUserWithEmailAddress(emailAddress string) (int64, error) {
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
