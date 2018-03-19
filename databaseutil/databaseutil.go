package databaseutil

import (
	"database/sql"
	// Notice that we’re loading the driver anonymously.
	// The driver registers itself as being available
	// to the database/sql package.
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var db *sql.DB

func Connect(dbUrl string) error {
	var err error

	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}

	// Quickly test if the connection to the database worked.
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func CreateUser(
	displayName string,
	emailAddress string,
	password string,
) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	sqlQuery := `
		INSERT INTO users (display_name, email_address, password, creation_time)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err = db.QueryRow(
		sqlQuery,
		displayName,
		emailAddress,
		hashedPassword,
		time.Now().UTC(),
	).Scan(&id)
	if err != nil {
		return -1, err
	}

	log.Printf("created new user with id '%d'", id)
	return id, nil
}

func AuthenticateUser(emailAddress string, password string) (bool, error) {
	sqlQuery := `
		SELECT password FROM users
		WHERE email_address = $1`

	var storedHashedPassword []byte
	err := db.QueryRow(sqlQuery, emailAddress).Scan(&storedHashedPassword)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword(
		storedHashedPassword,
		[]byte(password),
	); err != nil {
		return false, err
	}

	return true, nil
}
