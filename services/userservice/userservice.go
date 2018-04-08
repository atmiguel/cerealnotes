/*
Package userservice handles interactions of app with database layer.
*/
package userservice

import (
	"errors"
	"time"

	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
	"golang.org/x/crypto/bcrypt"
)

// EmailAddressAlreadyInUseError is returned when the email that was passed in
// cannot be used becaues it already being used by another user
var EmailAddressAlreadyInUseError = errors.New("Email address already in use")

// StoreNewUser takes a new user information and attempts to store it into the
// database
func StoreNewUser(
	displayName string,
	emailAddress *models.EmailAddress,
	password string,
) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	creationTime := time.Now().UTC()

	if err := databaseutil.InsertIntoUsersTable(
		displayName,
		emailAddress.String(),
		hashedPassword,
		creationTime,
	); err != nil {
		if err == databaseutil.UniqueConstraintError {
			return EmailAddressAlreadyInUseError
		}

		return err
	}

	return nil
}

// AuthenticateUserCredentials validates if the email address passwrod combo
// is valid. Returns nil on success and not nil depending on the error.
func AuthenticateUserCredentials(emailAddress *models.EmailAddress, password string) error {
	storedHashedPassword, err := databaseutil.GetPasswordForUserWithEmailAddress(
		emailAddress.String())
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(
		storedHashedPassword,
		[]byte(password),
	); err != nil {
		return err
	}

	return nil
}

// GetIdForUserWithEmailAddress returns the assoaited userId for a given
// email address.
func GetIdForUserWithEmailAddress(emailAddress *models.EmailAddress) (models.UserId, error) {
	userIdAsInt, err := databaseutil.GetIdForUserWithEmailAddress(emailAddress.String())
	if err != nil {
		return -1, err
	}

	return models.UserId(userIdAsInt), nil
}
