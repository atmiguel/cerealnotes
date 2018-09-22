package models

import (
	"encoding/json"
	"errors"
	"fmt"
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

type UserMap map[UserId]*User

func (userMap UserMap) ToJson() ([]byte, error) {
	// json doesn't support int indexed maps
	userByIdString := make(map[string]User, len(userMap))

	for id, user := range userMap {
		userByIdString[fmt.Sprint(id)] = *user
	}

	return json.Marshal(userByIdString)
}

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

	if _, err := db.execNoResults(sqlQuery, displayName, emailAddress.String(), hashedPassword, creationTime); err != nil {
		if err == UniqueConstraintError {
			return EmailAddressAlreadyInUseError
		}

		return err
	}

	return nil
}

func (db *DB) AuthenticateUserCredentials(emailAddress *EmailAddress, password string) error {
	sqlQuery := `
		SELECT password FROM app_user
		WHERE email_address = $1`

	var storedHashedPassword []byte

	if err := db.execOneResult(sqlQuery, &storedHashedPassword, emailAddress.String()); err != nil {
		if err == QueryResultContainedNoRowsError {
			return CredentialsNotAuthorizedError
		}

		return err
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

	var userId int64
	if err := db.execOneResult(sqlQuery, &userId, emailAddress.String()); err != nil {
		if err == QueryResultContainedNoRowsError {
			return 0, CredentialsNotAuthorizedError
		}
		return 0, err
	}

	return UserId(userId), nil
}

func (db *DB) GetAllUsersById() (UserMap, error) {
	sqlQuery := `
		SELECT id, display_name FROM app_user`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, convertPostgresError(err)
	}

	defer rows.Close()

	var userMap UserMap = make(map[UserId]*User)

	for rows.Next() {
		var tempId int64
		user := &User{}
		if err := rows.Scan(&tempId, &user.DisplayName); err != nil {
			return nil, convertPostgresError(err)
		}

		userMap[UserId(tempId)] = user

	}

	return userMap, nil

}
