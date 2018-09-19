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
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"

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
	equals(t, http.StatusOK, resp.StatusCode)
}

func TestAuthenticatedFlow(t *testing.T) {
	mockDb := &DiyMockDataStore{}
	env := &handlers.Environment{mockDb, []byte("")}

	server := httptest.NewServer(routers.DefineRoutes(env))
	defer server.Close()

	// Create testing client
	client := &http.Client{}
	{
		// jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
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

		ok(t, err)

		equals(t, http.StatusCreated, resp.StatusCode)
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
		ok(t, err)
		equals(t, http.StatusCreated, resp.StatusCode)

		type NoteResponse struct {
			NoteId int64 `json:"noteId"`
		}

		jsonNoteReponse := &NoteResponse{}

		err = json.NewDecoder(resp.Body).Decode(jsonNoteReponse)
		ok(t, err)

		equals(t, noteIdAsInt, jsonNoteReponse.NoteId)

		resp.Body.Close()
	}

	// Test get notes
	{
		mockDb.Func_GetMyUnpublishedNotes = func(userId models.UserId) (models.NoteMap, error) {
			if userIdAsInt != int64(userId) {
				return nil, errors.New("Invalid userId passed in")
			}

			return models.NoteMap(map[models.NoteId]*models.Note{
				models.NoteId(noteIdAsInt): &models.Note{
					AuthorId:     models.UserId(userIdAsInt),
					Content:      content,
					CreationTime: time.Now().UTC(),
				},
			}), nil

		}

		mockDb.Func_GetAllPublishedNotesVisibleBy = func(userId models.UserId) (map[int64]models.NoteMap, error) {
			if userIdAsInt != int64(userId) {
				return nil, errors.New("Invalid userId passed in")
			}

			return map[int64]models.NoteMap{
				1: models.NoteMap(map[models.NoteId]*models.Note{
					models.NoteId(44): &models.Note{
						AuthorId:     models.UserId(99),
						Content:      "another note",
						CreationTime: time.Now(),
					},
				}),
			}, nil

		}

		resp, err := client.Get(server.URL + paths.NoteApi)
		ok(t, err)
		equals(t, http.StatusOK, resp.StatusCode)
	}

	// Test Add category
	{
		type CategoryForm struct {
			NoteId   int64  `json:"noteId"`
			Category string `json:"category"`
		}

		metaCategory := models.META

		categoryForm := &CategoryForm{NoteId: noteIdAsInt, Category: metaCategory.String()}

		mockDb.Func_StoreNewNoteCategoryRelationship = func(noteId models.NoteId, cat models.Category) error {
			if int64(noteId) == noteIdAsInt && cat == metaCategory {
				return nil
			}

			return errors.New("Incorrect Data Arrived")
		}

		jsonValue, _ := json.Marshal(categoryForm)

		resp, err := client.Post(server.URL+paths.CategoryApi, "application/json", bytes.NewBuffer(jsonValue))
		ok(t, err)
		equals(t, http.StatusCreated, resp.StatusCode)
	}

	// Test publish notes
	{
		mockDb.Func_PublishNotes = func(userId models.UserId) error {
			return nil
		}
		// publish new api
		resp, err := client.Post(server.URL+paths.PublicationApi, "", nil)
		ok(t, err)
		equals(t, http.StatusCreated, resp.StatusCode)
	}

	// Test edit notes
	{
		type NoteUpdateForm struct {
			NoteId  int64  `json:"id"`
			Content string `json:"content"`
		}

		mockDb.Func_GetNoteById = func(models.NoteId) (*models.Note, error) {
			return &models.Note{
				AuthorId:     models.UserId(userIdAsInt),
				Content:      content,
				CreationTime: time.Now().UTC(),
			}, nil
		}

		mockDb.Func_UpdateNoteContent = func(models.NoteId, string) error {
			return nil
		}

		noteForm := &NoteUpdateForm{
			NoteId:  3,
			Content: "anything else",
		}

		jsonValue, _ := json.Marshal(noteForm)

		resp, err := sendPutRequest(client, server.URL+paths.NoteApi, "application/json", bytes.NewBuffer(jsonValue))

		ok(t, err)
		equals(t, http.StatusOK, resp.StatusCode)

	}

	// Delete note
	{
		mockDb.Func_GetUsersNotes = func(userId models.UserId) (models.NoteMap, error) {
			return models.NoteMap(map[models.NoteId]*models.Note{
				models.NoteId(noteIdAsInt): &models.Note{
					AuthorId:     models.UserId(userIdAsInt),
					Content:      content,
					CreationTime: time.Now(),
				},
			}), nil
		}

		mockDb.Func_DeleteNoteById = func(noteid models.NoteId) error {
			if int64(noteid) == noteIdAsInt {
				return nil
			}

			return errors.New("Somehow you didn't get the correct error")
		}

		resp, err := sendDeleteRequest(client, server.URL+paths.NoteApi+"?id="+strconv.FormatInt(noteIdAsInt, 10))
		ok(t, err)
		equals(t, http.StatusOK, resp.StatusCode)
	}

}

// func sendDeleteRequest(client *http.Client, myUrl string, contentType string, body io.Reader) (resp *http.Response, err error) {
func sendDeleteRequest(client *http.Client, myUrl string) (resp *http.Response, err error) {

	req, err := http.NewRequest("DELETE", myUrl, nil)

	if err != nil {
		return nil, err
	}

	return client.Do(req)
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
	Func_StoreNewNoteCategoryRelationship func(models.NoteId, models.Category) error
	Func_StoreNewUser                     func(string, *models.EmailAddress, string) error
	Func_AuthenticateUserCredentials      func(*models.EmailAddress, string) error
	Func_GetIdForUserWithEmailAddress     func(*models.EmailAddress) (models.UserId, error)
	Func_GetUsersNotes                    func(models.UserId) (models.NoteMap, error)
	Func_DeleteNoteById                   func(models.NoteId) error
	Func_GetMyUnpublishedNotes            func(models.UserId) (models.NoteMap, error)
	Func_GetAllUsersById                  func() (models.UserMap, error)
	Func_GetAllPublishedNotesVisibleBy    func(models.UserId) (map[int64]models.NoteMap, error)
	Func_PublishNotes                     func(models.UserId) error
	Func_StoreNewPublication              func(*models.Publication) (models.PublicationId, error)
	Func_GetNoteById                      func(models.NoteId) (*models.Note, error)
	Func_UpdateNoteContent                func(models.NoteId, string) error
	Func_UpdateNoteCategory               func(models.NoteId, models.Category) error
	Func_DeleteNoteCategory               func(models.NoteId) error
	Func_GetNoteCategory                  func(models.NoteId) (models.Category, error)
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

func (mock *DiyMockDataStore) GetUsersNotes(userId models.UserId) (models.NoteMap, error) {
	return mock.Func_GetUsersNotes(userId)
}

func (mock *DiyMockDataStore) DeleteNoteById(noteId models.NoteId) error {
	return mock.Func_DeleteNoteById(noteId)
}

func (mock *DiyMockDataStore) GetMyUnpublishedNotes(userId models.UserId) (models.NoteMap, error) {
	return mock.Func_GetMyUnpublishedNotes(userId)
}

func (mock *DiyMockDataStore) GetAllUsersById() (models.UserMap, error) {
	return mock.Func_GetAllUsersById()
}

func (mock *DiyMockDataStore) GetAllPublishedNotesVisibleBy(userId models.UserId) (map[int64]models.NoteMap, error) {
	return mock.Func_GetAllPublishedNotesVisibleBy(userId)
}

func (mock *DiyMockDataStore) PublishNotes(userId models.UserId) error {
	return mock.Func_PublishNotes(userId)
}

func (mock *DiyMockDataStore) StoreNewPublication(publication *models.Publication) (models.PublicationId, error) {
	return mock.Func_StoreNewPublication(publication)
}

func (mock *DiyMockDataStore) GetNoteById(noteId models.NoteId) (*models.Note, error) {
	return mock.Func_GetNoteById(noteId)
}

func (mock *DiyMockDataStore) UpdateNoteContent(noteId models.NoteId, content string) error {
	return mock.Func_UpdateNoteContent(noteId, content)
}

func (mock *DiyMockDataStore) UpdateNoteCategory(noteId models.NoteId, category models.Category) error {
	return mock.Func_UpdateNoteCategory(noteId, category)
}

func (mock *DiyMockDataStore) DeleteNoteCategory(noteId models.NoteId) error {
	return mock.Func_DeleteNoteCategory(noteId)
}

func (mock *DiyMockDataStore) GetNoteCategory(noteId models.NoteId) (models.Category, error) {
	return mock.Func_GetNoteCategory(noteId)
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
