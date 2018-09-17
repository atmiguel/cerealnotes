package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/paths"
	"github.com/atmiguel/cerealnotes/routers"
)

func TestLoginOrSignUpPage(t *testing.T) {
	mockDb := &DiyMockDataStore{}
	env := &handlers.Environment{mockDb, []byte("")}

	server := httptest.NewServer(routers.DefineRoutes(env))
	defer server.Close()

	resp, err := http.Get(server.URL)
	ok(t, err)

	// fmt.Println(ioutil.ReadAll(resp.Body))
	equals(t, http.StatusOK, resp.StatusCode)
}

func TestAuthenticatedFlow(t *testing.T) {
	mockDb := &DiyMockDataStore{}
	env := &handlers.Environment{mockDb, []byte("")}

	server := httptest.NewServer(routers.DefineRoutes(env))
	defer server.Close()

	// Create testing client
	// jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	jar, err := cookiejar.New(&cookiejar.Options{})

	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// Test login
	theEmail := "justsomeemail@gmail.com"
	thePassword := "worldsBestPassword"

	mockDb.Func_AuthenticateUserCredentials = func(email *models.EmailAddress, password string) error {
		if email.String() == theEmail && password == thePassword {
			return nil
		}

		return models.CredentialsNotAuthorizedError
	}

	mockDb.Func_GetIdForUserWithEmailAddress = func(email *models.EmailAddress) (models.UserId, error) {
		return models.UserId(1), nil
	}

	userValues := map[string]string{"emailAddress": theEmail, "password": thePassword}

	userJsonValue, _ := json.Marshal(userValues)

	resp, err := client.Post(server.URL+paths.SessionApi, "application/json", bytes.NewBuffer(userJsonValue))

	ok(t, err)

	equals(t, http.StatusCreated, resp.StatusCode)

	// Test Add Note
	noteValues := map[string]string{"content": "Dude I just said something cool"}
	noteIdAsInt := int64(33)

	mockDb.Func_StoreNewNote = func(*models.Note) (models.NoteId, error) {
		return models.NoteId(noteIdAsInt), nil
	}

	noteJsonValue, _ := json.Marshal(noteValues)

	resp, err = client.Post(server.URL+paths.NoteApi, "application/json", bytes.NewBuffer(noteJsonValue))
	ok(t, err)
	equals(t, http.StatusCreated, resp.StatusCode)

	type NoteResponse struct {
		NoteId int64 `json:"noteId"`
	}

	jsonNoteReponse := &NoteResponse{}

	err = json.NewDecoder(resp.Body).Decode(jsonNoteReponse)
	ok(t, err)

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))

	equals(t, noteIdAsInt, jsonNoteReponse.NoteId)

	resp.Body.Close()

}

// Helpers

type DiyMockDataStore struct {
	Func_StoreNewNote                     func(*models.Note) (models.NoteId, error)
	Func_StoreNewNoteCategoryRelationship func(models.NoteId, models.Category) error
	Func_StoreNewUser                     func(string, *models.EmailAddress, string) error
	Func_AuthenticateUserCredentials      func(*models.EmailAddress, string) error
	Func_GetIdForUserWithEmailAddress     func(*models.EmailAddress) (models.UserId, error)
}

func (mock *DiyMockDataStore) StoreNewNote(note *models.Note) (models.NoteId, error) {
	return mock.Func_StoreNewNote(note)
}

func (mock *DiyMockDataStore) StoreNewNoteCategoryRelationship(noteId models.NoteId, cat models.Category) error {
	return mock.Func_StoreNewNoteCategoryRelationship(noteId, cat)
}

func (mock *DiyMockDataStore) StoreNewUser(str1 string, email *models.EmailAddress, str2 string) error {
	return mock.Func_StoreNewUser(str1, email, str2)
}

func (mock *DiyMockDataStore) AuthenticateUserCredentials(email *models.EmailAddress, str string) error {
	return mock.Func_AuthenticateUserCredentials(email, str)
}

func (mock *DiyMockDataStore) GetIdForUserWithEmailAddress(email *models.EmailAddress) (models.UserId, error) {
	return mock.Func_GetIdForUserWithEmailAddress(email)
}

// assert fails the test if the condition is false.
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
