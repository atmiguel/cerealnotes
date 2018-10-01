package models_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/test_util"
)

var postgresUrl = "postgresql://localhost/cerealnotes_test?sslmode=disable"

const noteTable = "note"
const publicationTable = "publication"
const noteToPublicationTable = "note_to_publication_relationship"
const noteToCategoryTable = "note_to_category_relationship"
const userTable = "app_user"

var tables = []string{
	noteToPublicationTable,
	publicationTable,
	noteToCategoryTable,
	noteTable,
	userTable,
}

func ClearAllValuesInTable(db *models.DB) {
	for _, val := range tables {
		if err := ClearValuesInTable(db, val); err != nil {
			panic(err)
		}
	}
}

func ClearValuesInTable(db *models.DB, table string) error {
	// db.Query() doesn't allow variables to replace columns or table names.
	sqlQuery := fmt.Sprintf(`TRUNCATE %s CASCADE;`, table)

	_, err := db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func TestUser(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	test_util.Ok(t, err)
	ClearValuesInTable(db, userTable)

	displayName := "boby"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyOtherEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	test_util.Ok(t, err)

	_, err = db.GetIdForUserWithEmailAddress(emailAddress)
	test_util.Ok(t, err)

	err = db.AuthenticateUserCredentials(emailAddress, password)
	test_util.Ok(t, err)
}

func TestNote(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	test_util.Ok(t, err)
	ClearValuesInTable(db, userTable)
	ClearValuesInTable(db, noteTable)

	displayName := "bob"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	test_util.Ok(t, err)

	userId, err := db.GetIdForUserWithEmailAddress(emailAddress)
	test_util.Ok(t, err)

	note := &models.Note{AuthorId: userId, Content: "I'm a note", CreationTime: time.Now()}
	id, err := db.StoreNewNote(note)
	test_util.Ok(t, err)
	test_util.Assert(t, int64(id) > 0, "Note Id was not a valid index: "+strconv.Itoa(int(id)))
}

func TestCategory(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	test_util.Ok(t, err)
	ClearValuesInTable(db, userTable)
	ClearValuesInTable(db, noteTable)
	ClearValuesInTable(db, noteToCategoryTable)

	displayName := "bob"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisyetAnotherIsMyEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	test_util.Ok(t, err)

	userId, err := db.GetIdForUserWithEmailAddress(emailAddress)
	test_util.Ok(t, err)

	note := &models.Note{AuthorId: userId, Content: "I'm a note", CreationTime: time.Now()}
	noteId, err := db.StoreNewNote(note)
	test_util.Ok(t, err)

	err = db.StoreNewNoteCategoryRelationship(noteId, models.META)
	test_util.Ok(t, err)
}
