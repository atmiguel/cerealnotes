package models_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/atmiguel/cerealnotes/models"
)

var postgresUrl = "postgresql://localhost/test_db?sslmode=disable"

var tables = []string{
	"note_to_publication_relationship",
	"publication",
	"note_to_category_relationship",
	"note",
	"app_user",
}

func ClearAllValuesInTable(db *models.DB) {
	for _, val := range tables {
		if err := ClearValuesInTable(db, val); err != nil {
			panic(err)
		}
	}

}

func ClearValuesInTable(db *models.DB, table string) error {
	// db.Query() doesn't allow varaibles to replace columns or table names.
	sqlQuery := fmt.Sprintf(`TRUNCATE %s CASCADE;`, table)

	_, err := db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func TestUser(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)
	ClearAllValuesInTable(db)

	displayName := "boby"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyOtherEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	ok(t, err)

	_, err = db.GetIdForUserWithEmailAddress(emailAddress)
	ok(t, err)

	err = db.AuthenticateUserCredentials(emailAddress, password)
	ok(t, err)
}

func TestNote(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)
	ClearAllValuesInTable(db)

	displayName := "bob"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	ok(t, err)

	userId, err := db.GetIdForUserWithEmailAddress(emailAddress)
	ok(t, err)

	note := &models.Note{AuthorId: userId, Content: "I'm a note", CreationTime: time.Now()}
	id, err := db.StoreNewNote(note)
	ok(t, err)
	assert(t, int64(id) > 0, "Note Id was not a valid index: "+strconv.Itoa(int(id)))
}

func TestCategory(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)
	ClearAllValuesInTable(db)

	displayName := "bob"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisyetAnotherIsMyEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	ok(t, err)

	userId, err := db.GetIdForUserWithEmailAddress(emailAddress)
	ok(t, err)

	note := &models.Note{AuthorId: userId, Content: "I'm a note", CreationTime: time.Now()}
	noteId, err := db.StoreNewNote(note)
	ok(t, err)

	err = db.StoreNewNoteCategoryRelationship(noteId, models.META)
	ok(t, err)
}
