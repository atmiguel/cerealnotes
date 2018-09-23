package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/paths"
	"github.com/dgrijalva/jwt-go"
)

const oneWeek = time.Hour * 24 * 7
const credentialTimeoutDuration = oneWeek
const cerealNotesCookieName = "CerealNotesToken"
const baseTemplateName = "base"
const baseTemplateFile = "templates/base.tmpl"

var EmptyNoteContentError error = errors.New("Note content cannot be empty or just whitespace")
var NotYourNoteError error = errors.New("You are not the other of this note and therer for cannot preform this action")
var NoChangeError error = errors.New("The action you are trying to prefrom doesn't change anything")
var InvalidMethodError error = errors.New("This endpoint does not except that http method")

// JwtTokenClaim contains all claims required for authentication, including the standard JWT claims.
type JwtTokenClaim struct {
	models.UserId `json:"userId"`
	jwt.StandardClaims
}

type Environment struct {
	Db              models.Datastore
	TokenSigningKey []byte
}

type AuthenticatedRequestHandlerType func(
	*Environment,
	http.ResponseWriter,
	*http.Request,
	models.UserId,
) (error, int)

type UnauthenticatedEndpointHandlerType func(
	*Environment,
	http.ResponseWriter,
	*http.Request,
) (error, int)

// Wrappers
func AuthenticateOrRedirect(
	env *Environment,
	authenticatedHandlerFunc AuthenticatedRequestHandlerType,
	redirectPath string,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if userId, err := getUserIdFromJwtToken(env, request); err != nil {
			switch request.Method {
			// If not logged in, redirect to login page
			case http.MethodGet:
				http.Redirect(
					responseWriter,
					request,
					redirectPath,
					http.StatusTemporaryRedirect)
				return
			default:
				respondWithMethodNotAllowed(responseWriter, http.MethodGet)
			}
		} else {
			if err, errCode := authenticatedHandlerFunc(env, responseWriter, request, userId); err != nil {
				if errCode >= 500 {
					log.Print(err)
				}
				http.Error(responseWriter, err.Error(), errCode)
				return
			}
		}
	}
}

func AuthenticateOrReturnUnauthorized(
	env *Environment,
	authenticatedHandlerFunc AuthenticatedRequestHandlerType,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {

		if userId, err := getUserIdFromJwtToken(env, request); err != nil {
			responseWriter.Header().Set("WWW-Authenticate", `Bearer realm="`+request.URL.Path+`"`)
			http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
			return
		} else {
			if err, errCode := authenticatedHandlerFunc(env, responseWriter, request, userId); err != nil {
				if errCode >= 500 {
					log.Print(err)
				}
				http.Error(responseWriter, err.Error(), errCode)
				return
			}
		}
	}
}

func WrapUnauthenticatedEndpoint(env *Environment, handler UnauthenticatedEndpointHandlerType) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if err, errCode := handler(env, responseWriter, request); err != nil {
			if errCode >= 500 {
				log.Print(err)
			}
			http.Error(responseWriter, err.Error(), errCode)
			return
		}
	}
}

// UNAUTHENTICATED HANDLERS

// HandleLoginOrSignupPageRequest responds to unauthenticated GET requests with the login or signup page.
// For authenticated requests, it redirects to the home page.
func HandleLoginOrSignupPageRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
) (error, int) {
	switch request.Method {
	case http.MethodGet:
		if _, err := getUserIdFromJwtToken(env, request); err == nil {
			http.Redirect(
				responseWriter,
				request,
				paths.HomePage,
				http.StatusTemporaryRedirect)
			return nil, 0
		}

		parsedTemplate, err := template.ParseFiles(baseTemplateFile, "templates/login_or_signup.tmpl")
		if err != nil {
			return err, http.StatusInternalServerError
		}

		parsedTemplate.ExecuteTemplate(responseWriter, baseTemplateName, nil)

		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodGet)
	}
}

// API
func HandleUserApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
) (error, int) {
	type SignupForm struct {
		DisplayName  string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
		Password     string `json:"password"`
	}

	switch request.Method {
	case http.MethodPost:
		signupForm := new(SignupForm)

		if err := json.NewDecoder(request.Body).Decode(signupForm); err != nil {
			return err, http.StatusInternalServerError
		}

		var statusCode int
		if err := env.Db.StoreNewUser(
			signupForm.DisplayName,
			models.NewEmailAddress(signupForm.EmailAddress),
			signupForm.Password,
		); err != nil {
			if err == models.EmailAddressAlreadyInUseError {
				statusCode = http.StatusConflict
			} else {
				return err, http.StatusInternalServerError
			}
		} else {
			statusCode = http.StatusCreated
		}

		responseWriter.WriteHeader(statusCode)

		return nil, 0

	case http.MethodGet:

		if _, err := getUserIdFromJwtToken(env, request); err != nil {
			return err, http.StatusUnauthorized
		}

		usersById, err := env.Db.GetAllUsersById()
		if err != nil {
			return err, http.StatusInternalServerError
		}

		usersByIdJson, err := usersById.ToJson()
		if err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		fmt.Fprint(responseWriter, string(usersByIdJson))

		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodPost, http.MethodGet)
	}
}

// HandleSessionApiRequest responds to POST requests by authenticating and responding with a JWT.
// It responds to DELETE requests by expiring the client's cookie.
func HandleSessionApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
) (error, int) {
	type LoginForm struct {
		EmailAddress string `json:"emailAddress"`
		Password     string `json:"password"`
	}

	switch request.Method {
	case http.MethodPost:
		loginForm := new(LoginForm)

		if err := json.NewDecoder(request.Body).Decode(loginForm); err != nil {
			return err, http.StatusBadRequest
		}

		if err := env.Db.AuthenticateUserCredentials(
			models.NewEmailAddress(loginForm.EmailAddress),
			loginForm.Password,
		); err != nil {
			statusCode := http.StatusInternalServerError
			if err == models.CredentialsNotAuthorizedError {
				statusCode = http.StatusUnauthorized
			}
			return err, statusCode
		}

		// Set our cookie to have a valid JWT Token as the value
		{
			userId, err := env.Db.GetIdForUserWithEmailAddress(models.NewEmailAddress(loginForm.EmailAddress))
			if err != nil {
				return err, http.StatusInternalServerError
			}

			token, err := CreateTokenAsString(env, userId, credentialTimeoutDuration)
			if err != nil {
				return err, http.StatusInternalServerError
			}

			expirationTime := time.Now().Add(credentialTimeoutDuration)

			cookie := http.Cookie{
				Name:     cerealNotesCookieName,
				Value:    token,
				Path:     "/",
				Expires:  expirationTime,
				HttpOnly: true,
			}

			http.SetCookie(responseWriter, &cookie)
		}

		responseWriter.WriteHeader(http.StatusCreated)

		return nil, 0

	case http.MethodDelete:
		// Cookie will overwrite existing cookie then delete itself
		cookie := http.Cookie{
			Name:     cerealNotesCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		}

		http.SetCookie(responseWriter, &cookie)
		responseWriter.WriteHeader(http.StatusOK)
		fmt.Fprint(responseWriter, "user successfully logged out")

		return nil, 0

	default:
		return respondWithMethodNotAllowed(
			responseWriter,
			http.MethodPost,
			http.MethodDelete)
	}
}

func HandlePublicationApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) (error, int) {
	switch request.Method {
	case http.MethodPost:
		if err := env.Db.PublishNotes(userId); err != nil {
			return err, http.StatusInternalServerError
		}
		responseWriter.WriteHeader(http.StatusCreated)

		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodPost)
	}
}

func HandleNoteApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) (error, int) {

	type NoteForm struct {
		Content string `json:"content"`
	}
	switch request.Method {

	case http.MethodGet:

		publishedNotes, err := env.Db.GetAllPublishedNotesVisibleBy(userId)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		myUnpublishedNotes, err := env.Db.GetMyUnpublishedNotes(userId)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		// fmt.Println("number of published notes")
		// fmt.Println(len(publishedNotes))
		// fmt.Println("number of unpublished notes")
		// fmt.Println(len(myUnpublishedNotes))

		allNotes := myUnpublishedNotes

		// TODO figure out how to surface the publication number

		// for publicationNumber, noteMap := range publishedNotes {
		for _, noteMap := range publishedNotes {
			for id, note := range noteMap {
				allNotes[id] = note
			}
		}

		notesInJson, err := allNotes.ToJson()
		if err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)

		fmt.Fprint(responseWriter, string(notesInJson))

		return nil, 0

	case http.MethodPost:

		noteForm := new(NoteForm)

		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			return err, http.StatusBadRequest
		}

		if len(strings.TrimSpace(noteForm.Content)) == 0 {
			return EmptyNoteContentError, http.StatusBadRequest
		}

		note := &models.Note{
			AuthorId:     models.UserId(userId),
			Content:      strings.TrimSpace(noteForm.Content),
			CreationTime: time.Now().UTC(),
		}

		noteId, err := env.Db.StoreNewNote(note)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		type NoteResponse struct {
			NoteId int64 `json:"noteId"`
		}

		noteString, err := json.Marshal(&NoteResponse{NoteId: int64(noteId)})
		if err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusCreated)

		fmt.Fprint(responseWriter, string(noteString))

		return nil, 0

	case http.MethodPut:

		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		noteId := models.NoteId(id)

		noteForm := new(NoteForm)
		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			return err, http.StatusBadRequest
		}

		if noteId < 1 {
			return models.NoNoteFoundError, http.StatusBadRequest
		}

		note, err := env.Db.GetNoteById(noteId)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		if note.AuthorId != userId {
			return NotYourNoteError, http.StatusUnauthorized
		}

		content := strings.TrimSpace(noteForm.Content)
		if len(content) == 0 {
			return EmptyNoteContentError, http.StatusBadRequest
		}

		if content == note.Content {
			return NoChangeError, http.StatusBadRequest
		}

		if err := env.Db.UpdateNoteContent(noteId, content); err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.WriteHeader(http.StatusOK)

		return nil, 0

	case http.MethodDelete:

		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		if err != nil {
			return err, http.StatusBadRequest
		}

		noteId := models.NoteId(id)

		noteMap, err := env.Db.GetUsersNotes(userId)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		if _, ok := noteMap[noteId]; !ok {
			return models.NoNoteFoundError, http.StatusInternalServerError
		}

		err = env.Db.DeleteNoteById(noteId)
		if err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.WriteHeader(http.StatusOK)

		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodGet, http.MethodPost, http.MethodDelete)
	}
}

func HandleCategoryApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) (error, int) {
	switch request.Method {
	case http.MethodGet:

		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		noteId := models.NoteId(id)

		var categoryString string

		category, err := env.Db.GetNoteCategory(noteId)
		if err != nil {
			if err == models.QueryResultContainedNoRowsError {
				// you are trying to a a non set category return the empty string
				categoryString = ""

			} else {
				return err, http.StatusInternalServerError
			}
		} else {
			categoryString = category.String()
		}

		type categoryObj struct {
			Category string `json:"category"`
		}

		jsonValue, err := json.Marshal(&categoryObj{Category: categoryString})
		if err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)

		fmt.Fprint(responseWriter, string(jsonValue))

		return nil, 0

	case http.MethodPost:

		type CategoryForm struct {
			NoteId   int64  `json:"noteId"`
			Category string `json:"category"`
		}

		noteForm := new(CategoryForm)

		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			return err, http.StatusBadRequest
		}

		category, err := models.DeserializeCategory(strings.ToLower(noteForm.Category))

		if err != nil {
			return err, http.StatusBadRequest
		}

		if err := env.Db.StoreNewNoteCategoryRelationship(models.NoteId(noteForm.NoteId), category); err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.WriteHeader(http.StatusCreated)

		return nil, 0

	case http.MethodPut:
		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		noteId := models.NoteId(id)

		_, err = env.Db.GetNoteById(noteId)
		if err != nil {
			if err == models.NoNoteFoundError {
				return err, http.StatusBadRequest
			}
			return err, http.StatusInternalServerError
		}

		type CategoryForm struct {
			Category string `json:"category"`
		}

		noteForm := new(CategoryForm)

		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			return err, http.StatusBadRequest
		}

		category, err := models.DeserializeCategory(strings.ToLower(noteForm.Category))
		if err != nil {
			return err, http.StatusBadRequest
		}

		if err := env.Db.UpdateNoteCategory(noteId, category); err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.WriteHeader(http.StatusOK)

		return nil, 0

	case http.MethodDelete:

		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		if err != nil {
			return err, http.StatusBadRequest
		}

		noteId := models.NoteId(id)

		if err := env.Db.DeleteNoteCategory(noteId); err != nil {
			return err, http.StatusInternalServerError
		}

		responseWriter.WriteHeader(http.StatusOK)

		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodPost, http.MethodPut, http.MethodDelete)
	}

}

func RedirectToPathHandler(
	path string,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			http.Redirect(
				responseWriter,
				request,
				path,
				http.StatusTemporaryRedirect)
			return
		default:
			err, errCode := respondWithMethodNotAllowed(responseWriter, http.MethodGet)
			if errCode >= 500 {
				log.Print(err)
			}
			http.Error(responseWriter, err.Error(), errCode)
			return
		}
	}
}

// AUTHENTICATED HANDLERS

func HandleHomePageRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) (error, int) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles(baseTemplateFile, "templates/home.tmpl")
		if err != nil {
			return err, http.StatusInternalServerError
		}

		parsedTemplate.ExecuteTemplate(responseWriter, baseTemplateName, userId)

		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodGet)
	}
}

func HandleNotesPageRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) (error, int) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles(baseTemplateFile, "templates/notes.tmpl")
		if err != nil {
			return err, http.StatusInternalServerError
		}

		parsedTemplate.ExecuteTemplate(responseWriter, baseTemplateName, userId)
		return nil, 0

	default:
		return respondWithMethodNotAllowed(responseWriter, http.MethodGet)
	}
}

// PRIVATE

func respondWithMethodNotAllowed(
	responseWriter http.ResponseWriter,
	allowedMethod string,
	otherAllowedMethods ...string,
) (error, int) {
	allowedMethods := append([]string{allowedMethod}, otherAllowedMethods...)
	allowedMethodsString := strings.Join(allowedMethods, ", ")

	responseWriter.Header().Set("Allow", allowedMethodsString)

	return InvalidMethodError, http.StatusMethodNotAllowed
}
