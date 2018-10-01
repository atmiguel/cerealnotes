package main_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/paths"
	"github.com/atmiguel/cerealnotes/routers"
	"github.com/atmiguel/cerealnotes/test_util"
)

func TestLoginOrSignUpPage(t *testing.T) {
	mockDb := &DiyMockDataStore{}
	env := &handlers.Environment{mockDb, []byte("")}

	server := httptest.NewServer(routers.DefineRoutes(env))
	defer server.Close()

	resp, err := http.Get(server.URL)
	test_util.Ok(t, err)
	test_util.Equals(t, http.StatusOK, resp.StatusCode)
}

func TestAuthenticatedFlow(t *testing.T) {
	mockDb := &DiyMockDataStore{}
	env := &handlers.Environment{mockDb, []byte("")}

	server := httptest.NewServer(routers.DefineRoutes(env))
	defer server.Close()

	// Create testing client
	client := &http.Client{}
	{
		jar, err := cookiejar.New(&cookiejar.Options{})

		if err != nil {
			panic(err)
		}

		client.Jar = jar
	}

	// Test login
	userIdAsInt := int64(1)
	{
		theEmail := "justsomeemail@gmail.com"
		thePassword := "worldsBestPassword"

		mockDb.Func_AuthenticateUserCredentials = func(email *models.EmailAddress, password string) error {
			if email.String() == theEmail && password == thePassword {
				return nil
			}

			return models.CredentialsNotAuthorizedError
		}

		mockDb.Func_GetIdForUserWithEmailAddress = func(email *models.EmailAddress) (models.UserId, error) {
			return models.UserId(userIdAsInt), nil
		}

		userValues := map[string]string{"emailAddress": theEmail, "password": thePassword}

		userJsonValue, _ := json.Marshal(userValues)

		resp, err := client.Post(server.URL+paths.SessionApi, "application/json", bytes.NewBuffer(userJsonValue))

		test_util.Ok(t, err)

		test_util.Equals(t, http.StatusCreated, resp.StatusCode)
	}

	// Test Add Note
	noteIdAsInt := int64(33)
	content := "Duuude I just said something cool"
	{
		noteValues := map[string]string{"content": content}

		mockDb.Func_StoreNewNote = func(*models.Note) (models.NoteId, error) {
			return models.NoteId(noteIdAsInt), nil
		}

		noteJsonValue, _ := json.Marshal(noteValues)

		resp, err := client.Post(server.URL+paths.NoteApi, "application/json", bytes.NewBuffer(noteJsonValue))
		test_util.Ok(t, err)
		test_util.Equals(t, http.StatusCreated, resp.StatusCode)

		type NoteResponse struct {
			NoteId int64 `json:"noteId"`
		}

		jsonNoteReponse := &NoteResponse{}

		err = json.NewDecoder(resp.Body).Decode(jsonNoteReponse)
		test_util.Ok(t, err)

		test_util.Equals(t, noteIdAsInt, jsonNoteReponse.NoteId)

		resp.Body.Close()
	}

	// Test get notes
	{
		resp, err := client.Get(server.URL + paths.NoteApi)
		test_util.Ok(t, err)
		test_util.Equals(t, http.StatusOK, resp.StatusCode)

		// TODO when we implement a real get notes feature we should enhance this code.
	}

	// Test Add category
	{
		type NoteCategoryForm struct {
			NoteCategory string `json:"category"`
		}

		metaNoteCategory := models.META

		categoryForm := &NoteCategoryForm{NoteCategory: metaNoteCategory.String()}

		mockDb.Func_StoreNewNoteCategoryRelationship = func(noteId models.NoteId, cat models.NoteCategory) error {
			if int64(noteId) == noteIdAsInt && cat == metaNoteCategory {
				return nil
			}

			return errors.New("Incorrect Data Arrived")
		}

		jsonValue, _ := json.Marshal(categoryForm)

		resp, err := sendPutRequest(client, server.URL+paths.NoteCategoryApi+"?id="+strconv.FormatInt(noteIdAsInt, 10), "application/json", bytes.NewBuffer(jsonValue))
		test_util.Ok(t, err)
		test_util.Equals(t, http.StatusCreated, resp.StatusCode)
	}

}

func sendPutRequest(client *http.Client, myUrl string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", myUrl, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	return client.Do(req)
}

func printBody(resp *http.Response) {
	buf, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		fmt.Print("bodyErr ", bodyErr.Error())
		return
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	fmt.Printf("BODY: %q", rdr1)
	resp.Body = rdr2
}

// Helpers

type DiyMockDataStore struct {
	Func_StoreNewNote                     func(*models.Note) (models.NoteId, error)
	Func_StoreNewNoteCategoryRelationship func(models.NoteId, models.NoteCategory) error
	Func_StoreNewUser                     func(string, *models.EmailAddress, string) error
	Func_AuthenticateUserCredentials      func(*models.EmailAddress, string) error
	Func_GetIdForUserWithEmailAddress     func(*models.EmailAddress) (models.UserId, error)
}

func (mock *DiyMockDataStore) StoreNewNote(note *models.Note) (models.NoteId, error) {
	return mock.Func_StoreNewNote(note)
}

func (mock *DiyMockDataStore) StoreNewNoteCategoryRelationship(noteId models.NoteId, cat models.NoteCategory) error {
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
