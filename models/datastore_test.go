package models_test

import (
	"testing"
	"time"
	"strconv"

	"github.com/atmiguel/cerealnotes/models"
)

var postgresUrl = "postgresql://localhost/test_db?sslmode=disable"

func ClearAllValuesInTable(*models.DB) {
	// todo call trucate_tables.sql
}

func TestUser(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)

	displayName := "boby"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyOtherEmail@gmail.com")

	err = db.StoreNewUser(displayName,emailAddress,password)
	ok(t, err) 

	_, err = db.GetIdForUserWithEmailAddress(emailAddress)
	ok(t, err)

	err = db.AuthenticateUserCredentials(emailAddress, password)
	ok(t, err)
}

func TestNote(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)

	displayName := "bob"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyEmail@gmail.com")

	err = db.StoreNewUser(displayName,emailAddress,password)
	ok(t, err) 

	userId, err := db.GetIdForUserWithEmailAddress(emailAddress)
	ok(t, err)

	note := &models.Note{AuthorId: userId, Content: "I'm a note", CreationTime: time.Now()}
	id, err := db.StoreNewNote(note)
	ok(t, err)
	assert(t, int64(id) > 0, "Note Id was not a valid index: "+strconv.Itoa(int(id)))
}
