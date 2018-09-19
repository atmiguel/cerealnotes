package models_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/atmiguel/cerealnotes/models"
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
	ClearValuesInTable(db, userTable)

	displayName := "boby"
	password := "aPassword"
	emailAddress := models.NewEmailAddress("thisIsMyOtherEmail@gmail.com")

	err = db.StoreNewUser(displayName, emailAddress, password)
	ok(t, err)

	id, err := db.GetIdForUserWithEmailAddress(emailAddress)
	ok(t, err)

	err = db.AuthenticateUserCredentials(emailAddress, password)
	ok(t, err)

	userMap, err := db.GetAllUsersById()
	ok(t, err)

	equals(t, 1, len(userMap))

	user, isOk := userMap[id]
	assert(t, isOk, "Expected UserId missing")

	equals(t, displayName, user.DisplayName)
}

func TestNote(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)
	ClearValuesInTable(db, userTable)
	ClearValuesInTable(db, noteTable)

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

	notemap, err := db.GetMyUnpublishedNotes(userId)
	ok(t, err)

	retrievedNote, isOk := notemap[id]
	assert(t, isOk, "Expected NoteId missing")

	equals(t, note.AuthorId, retrievedNote.AuthorId)
	equals(t, note.Content, retrievedNote.Content)

	err = db.DeleteNoteById(id)
	ok(t, err)
}

func TestPublication(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)
	ClearValuesInTable(db, userTable)
	ClearValuesInTable(db, noteTable)
	ClearValuesInTable(db, publicationTable)
	ClearValuesInTable(db, noteToPublicationTable)

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

	publicationToNoteMap, err := db.GetAllPublishedNotesVisibleBy(userId)
	ok(t, err)

	equals(t, 0, len(publicationToNoteMap))

	err = db.PublishNotes(userId)
	ok(t, err)

	publicationToNoteMap, err = db.GetAllPublishedNotesVisibleBy(userId)
	equals(t, 1, len(publicationToNoteMap))
}

func TestCategory(t *testing.T) {
	db, err := models.ConnectToDatabase(postgresUrl)
	ok(t, err)
	ClearValuesInTable(db, userTable)
	ClearValuesInTable(db, noteTable)
	ClearValuesInTable(db, noteToCategoryTable)

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

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
