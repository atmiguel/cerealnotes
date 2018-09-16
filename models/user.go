package models

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserId int64

type User struct {
	DisplayName string `json:"displayName"`
}

// EmailAddress ensures that email addresses are always formatted properly within the backend.
type EmailAddress struct {
	emailAddressAsString string
}

func NewEmailAddress(emailAddressAsString string) *EmailAddress {
	return &EmailAddress{emailAddressAsString: strings.ToLower(emailAddressAsString)}
}

func (emailAddress *EmailAddress) String() string {
	return emailAddress.emailAddressAsString
}

var EmailAddressAlreadyInUseError = errors.New("Email address already in use")

var CredentialsNotAuthorizedError = errors.New("The provided credentials were not found")

//

func (db *DB) StoreNewUser(
	displayName string,
	emailAddress *EmailAddress,
	password string,
) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	creationTime := time.Now().UTC()

	sqlQuery := `
		INSERT INTO app_user (display_name, email_address, password, creation_time)
		VALUES ($1, $2, $3, $4)`

	rows, err := db.Query(sqlQuery, displayName, emailAddress.String(), hashedPassword, creationTime)
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return convertPostgresError(err)
	}

	return nil
}

func (db *DB) AuthenticateUserCredentials(emailAddress *EmailAddress, password string) error {
	sqlQuery := `
		SELECT password FROM app_user
		WHERE email_address = $1`

	rows, err := db.Query(sqlQuery, emailAddress.String())
	if err != nil {
		return convertPostgresError(err)
	}
	defer rows.Close()

	var storedHashedPassword []byte
	for rows.Next() {
		if storedHashedPassword != nil {
			return QueryResultContainedMultipleRowsError
		}

		if err := rows.Scan(&storedHashedPassword); err != nil {
			return err
		}
	}

	if storedHashedPassword == nil {
		return QueryResultContainedNoRowsError
	}

	if err := bcrypt.CompareHashAndPassword(
		storedHashedPassword,
		[]byte(password),
	); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return CredentialsNotAuthorizedError
		}

		return err
	}

	return nil
}

func (db *DB) GetIdForUserWithEmailAddress(emailAddress *EmailAddress) (UserId, error) {
	sqlQuery := `
		SELECT id FROM app_user
		WHERE email_address = $1`

	rows, err := db.Query(sqlQuery, emailAddress.String())
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

	return UserId(userId), nil
}
