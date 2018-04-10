/*
Package userservice handles interactions with database layer.
*/
package userservice

import (
	"errors"
	"time"

	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
	"golang.org/x/crypto/bcrypt"
)

var EmailAddressAlreadyInUseError = errors.New("Email address already in use")

var CredentialsNotAuthorizedError = errors.New("The provided credentials were not found")

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
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return CredentialsNotAuthorizedError
		}

		return err
	}

	return nil
}

func GetIdForUserWithEmailAddress(emailAddress *models.EmailAddress) (models.UserId, error) {
	userIdAsInt, err := databaseutil.GetIdForUserWithEmailAddress(emailAddress.String())
	if err != nil {
		return -1, err
	}

	return models.UserId(userIdAsInt), nil
}
