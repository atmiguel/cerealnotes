package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
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

// JwtTokenClaim contains all claims required for authentication, including the standard JWT claims.
type JwtTokenClaim struct {
	models.UserId `json:"userId"`
	jwt.StandardClaims
}

type Environment struct {
	Db              models.Datastore
	TokenSigningKey []byte
}

func WrapUnauthenticatedEndpoint(env *Environment, handler UnauthenticatedEndpointHandlerType) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		handler(env, responseWriter, request)
	}
}

// UNAUTHENTICATED HANDLERS

// HandleLoginOrSignupPageRequest responds to unauthenticated GET requests with the login or signup page.
// For authenticated requests, it redirects to the home page.
func HandleLoginOrSignupPageRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
		if _, err := getUserIdFromJwtToken(env, request); err == nil {
			http.Redirect(
				responseWriter,
				request,
				paths.HomePage,
				http.StatusTemporaryRedirect)
			return
		}

		parsedTemplate, err := template.ParseFiles(baseTemplateFile, "templates/login_or_signup.tmpl")
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTemplate.ExecuteTemplate(responseWriter, baseTemplateName, nil)

	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodGet)
	}
}

func HandleUserApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	type SignupForm struct {
		DisplayName  string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
		Password     string `json:"password"`
	}

	switch request.Method {
	case http.MethodPost:
		signupForm := new(SignupForm)

		if err := json.NewDecoder(request.Body).Decode(signupForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
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
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			statusCode = http.StatusCreated
		}

		responseWriter.WriteHeader(statusCode)

	case http.MethodGet:

		if _, err := getUserIdFromJwtToken(env, request); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
			return
		}

		usersById, err := env.Db.GetAllUsersById()
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		usersByIdJson, err := usersById.ToJson()
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		fmt.Fprint(responseWriter, string(usersByIdJson))

	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodPost, http.MethodGet)
	}
}

// HandleSessionApiRequest responds to POST requests by authenticating and responding with a JWT.
// It responds to DELETE requests by expiring the client's cookie.
func HandleSessionApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	type LoginForm struct {
		EmailAddress string `json:"emailAddress"`
		Password     string `json:"password"`
	}

	switch request.Method {
	case http.MethodPost:
		loginForm := new(LoginForm)

		if err := json.NewDecoder(request.Body).Decode(loginForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := env.Db.AuthenticateUserCredentials(
			models.NewEmailAddress(loginForm.EmailAddress),
			loginForm.Password,
		); err != nil {
			statusCode := http.StatusInternalServerError
			if err == models.CredentialsNotAuthorizedError {
				statusCode = http.StatusUnauthorized
			}
			http.Error(responseWriter, err.Error(), statusCode)
			return
		}

		// Set our cookie to have a valid JWT Token as the value
		{
			userId, err := env.Db.GetIdForUserWithEmailAddress(models.NewEmailAddress(loginForm.EmailAddress))
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			token, err := CreateTokenAsString(env, userId, credentialTimeoutDuration)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
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

	default:
		respondWithMethodNotAllowed(
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
) {
	switch request.Method {
	case http.MethodPost:
		if err := env.Db.PublishNotes(userId); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		responseWriter.WriteHeader(http.StatusCreated)

	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodPost)
	}
}

func HandleNoteApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) {
	switch request.Method {
	case http.MethodGet:

		publishedNotes, err := env.Db.GetAllPublishedNotesVisibleBy(userId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		myUnpublishedNotes, err := env.Db.GetMyUnpublishedNotes(userId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("number of published notes")
		fmt.Println(len(publishedNotes))
		fmt.Println("number of unpublished notes")
		fmt.Println(len(myUnpublishedNotes))

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
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)

		fmt.Fprint(responseWriter, string(notesInJson))

	case http.MethodPost:
		type NoteForm struct {
			Content string `json:"content"`
		}

		noteForm := new(NoteForm)

		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		if len(strings.TrimSpace(noteForm.Content)) == 0 {
			http.Error(responseWriter, "Note content cannot be empty or just whitespace", http.StatusBadRequest)
			return
		}

		note := &models.Note{
			AuthorId:     models.UserId(userId),
			Content:      strings.TrimSpace(noteForm.Content),
			CreationTime: time.Now().UTC(),
		}

		noteId, err := env.Db.StoreNewNote(note)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		type NoteResponse struct {
			NoteId int64 `json:"noteId"`
		}

		noteString, err := json.Marshal(&NoteResponse{NoteId: int64(noteId)})
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusCreated)

		fmt.Fprint(responseWriter, string(noteString))

	case http.MethodPut:
		type NoteForm struct {
			Id      int64  `json:"id"`
			Content string `json:"content"`
		}

		noteForm := new(NoteForm)
		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		if noteForm.Id < 1 {
			http.Error(responseWriter, "Invalid Note Id", http.StatusBadRequest)
			return
		}

		noteId := models.NoteId(noteForm.Id)
		note, err := env.Db.GetNoteById(noteId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		if note.AuthorId != userId {
			http.Error(responseWriter, "You can only edit notes of which you are the author", http.StatusUnauthorized)
			return
		}

		content := strings.TrimSpace(noteForm.Content)
		if len(content) == 0 {
			http.Error(responseWriter, "Note content cannot be empty or just whitespace", http.StatusBadRequest)
			return
		}

		if content == note.Content {
			http.Error(responseWriter, "Note content is the same as existing content", http.StatusBadRequest)
			return
		}

		if err := env.Db.UpdateNoteContent(noteId, content); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
	case http.MethodDelete:

		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		noteId := models.NoteId(id)

		noteMap, err := env.Db.GetUsersNotes(userId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := noteMap[noteId]; !ok {
			errorString := "No note with that Id written by you was found"
			http.Error(responseWriter, errorString, http.StatusBadRequest)
			return
		}

		err = env.Db.DeleteNoteById(noteId)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)

	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodGet, http.MethodPost, http.MethodDelete)
	}
}

func HandleCategoryApiRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) {
	switch request.Method {
	case http.MethodPost:

		type CategoryForm struct {
			NoteId   int64  `json:"noteId"`
			Category string `json:"category"`
		}

		noteForm := new(CategoryForm)

		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		category, err := models.DeserializeCategory(strings.ToLower(noteForm.Category))

		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		if err := env.Db.StoreNewNoteCategoryRelationship(models.NoteId(noteForm.NoteId), category); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		type CategoryForm struct {
			NoteId   int64  `json:"noteId"`
			Category string `json:"category"`
		}

		noteForm := new(CategoryForm)

		if err := json.NewDecoder(request.Body).Decode(noteForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		category, err := models.DeserializeCategory(strings.ToLower(noteForm.Category))

		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		if err := env.Db.UpdateNoteCategory(models.NoteId(noteForm.NoteId), category); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusCreated)

	case http.MethodDelete:

		type DeleteForm struct {
			NoteId int64 `json:"noteId"`
		}

		deleteForm := new(DeleteForm)

		if err := json.NewDecoder(request.Body).Decode(deleteForm); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		if err := env.Db.DeleteNoteCategory(models.NoteId(deleteForm.NoteId)); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusCreated)

	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodPost, http.MethodPut, http.MethodDelete)
	}

}

type AuthenticatedRequestHandlerType func(
	*Environment,
	http.ResponseWriter,
	*http.Request,
	models.UserId,
)

type UnauthenticatedEndpointHandlerType func(
	*Environment,
	http.ResponseWriter,
	*http.Request,
)

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
			authenticatedHandlerFunc(env, responseWriter, request, userId)
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
		} else {
			authenticatedHandlerFunc(env, responseWriter, request, userId)
		}
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
			respondWithMethodNotAllowed(responseWriter, http.MethodGet)
		}
	}
}

// AUTHENTICATED HANDLERS

func HandleHomePageRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles(baseTemplateFile, "templates/home.tmpl")
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTemplate.ExecuteTemplate(responseWriter, baseTemplateName, userId)
	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodGet)
	}
}

func HandleNotesPageRequest(
	env *Environment,
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles(baseTemplateFile, "templates/notes.tmpl")
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTemplate.ExecuteTemplate(responseWriter, baseTemplateName, userId)

	default:
		respondWithMethodNotAllowed(responseWriter, http.MethodGet)
	}
}

// PRIVATE

func respondWithMethodNotAllowed(
	responseWriter http.ResponseWriter,
	allowedMethod string,
	otherAllowedMethods ...string,
) {
	allowedMethods := append([]string{allowedMethod}, otherAllowedMethods...)
	allowedMethodsString := strings.Join(allowedMethods, ", ")

	responseWriter.Header().Set("Allow", allowedMethodsString)
	statusCode := http.StatusMethodNotAllowed

	http.Error(responseWriter, http.StatusText(statusCode), statusCode)
}
