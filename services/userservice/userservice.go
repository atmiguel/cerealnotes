package userservice

import (
	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func StoreNewUser(
	displayName string,
	emailAddress string,
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
		emailAddress,
		hashedPassword,
		creationTime,
	); err != nil {
		return err
	}

	return nil
}

func AuthenticateUserCredentials(emailAddress string, password string) error {
	storedHashedPassword, err := databaseutil.GetPasswordForUserWithEmailAddress(
		emailAddress)
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

func GetIdForUserWithEmailAddress(emailAddress string) (models.UserId, error) {
	userIdAsInt, err := databaseutil.GetIdForUserWithEmailAddress(emailAddress)
	if err != nil {
		return -1, err
	}

	return models.UserId(userIdAsInt), nil
}
