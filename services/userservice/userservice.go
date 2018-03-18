package userservice

import (
	"github.com/atmiguel/cerealnotes/databaseutil"
	"github.com/atmiguel/cerealnotes/models/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func CreateUser(
	displayName string,
	emailAddress string,
	password string,
) (user.UserId, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	creationTime := time.Now().UTC()

	return databaseutil.InsertIntoUsersTable(
		displayName,
		emailAddress,
		hashedPassword,
		creationTime)
}

func AuthenticateUser(emailAddress string, password string) error {
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
