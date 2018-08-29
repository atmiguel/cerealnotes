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

	if err := databaseutil.InsertIntoUserTable(
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

func GetUsersById() (map[models.UserId]*models.User, error) {
	userData, err := databaseutil.GetAllUserData()
	if err != nil {
		return nil, err
	}

	usersById := make(map[models.UserId]*models.User)
	for _, userDatum := range userData {
		usersById[models.UserId(userDatum.Id)] = &models.User{DisplayName: userDatum.DisplayName}
	}

	return usersById, nil
}

func AuthenticateUserCredentials(emailAddress *models.EmailAddress, password string) error {
	storedHashedPassword, err := databaseutil.GetPasswordForUserWithEmailAddress(emailAddress.String())
	if err != nil {
		if err == databaseutil.QueryResultContainedMultipleRowsError {
			return err // would normally throw a runtime here
		}

		if err == databaseutil.QueryResultContainedNoRowsError {
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

func GetIdForUserWithEmailAddress(emailAddress *models.EmailAddress) (models.UserId, error) {
	userIdAsInt, err := databaseutil.GetIdForUserWithEmailAddress(emailAddress.String())
	if err != nil {
		if err == databaseutil.QueryResultContainedMultipleRowsError {
			return 0, err // would normally throw a runtime here
		}

		if err == databaseutil.QueryResultContainedNoRowsError {
			return 0, err
		}

		return 0, err
	}

	return models.UserId(userIdAsInt), nil
}
