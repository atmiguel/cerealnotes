package databaseutil

import (
	"database/sql"
	"github.com/atmiguel/cerealnotes/models/user"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var db *sql.DB

func init() {
	var databaseUrl string
	{
		environmentVariableName := "DATABASE_URL"
		databaseUrl = os.Getenv(environmentVariableName)

		if len(databaseUrl) == 0 {
			log.Fatalf("environment variable %s not set", environmentVariableName)
		}
	}

	{
		tempDb, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			log.Fatal(err)
		}

		db = tempDb
	}

	// Quickly test if the connection to the database worked.
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func InsertIntoUsersTable(
	displayName string,
	emailAddress string,
	password []byte,
	creationTime time.Time,
) (user.UserId, error) {
	var row *sql.Row
	{
		sqlQuery := `
			INSERT INTO users (display_name, email_address, password, creation_time)
			VALUES ($1, $2, $3, $4) RETURNING id `

		row = db.QueryRow(sqlQuery, displayName, emailAddress, password, creationTime)
	}

	var userId user.UserId
	if err := row.Scan(&userId); err != nil {
		// TODO handle err.ErrNoRows
		return -1, err
	}

	return userId, nil
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
		// TODO handle err.ErrNoRows
		return nil, err
	}

	return password, nil
}
