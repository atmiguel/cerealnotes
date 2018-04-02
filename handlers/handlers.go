package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/paths"
	"github.com/atmiguel/cerealnotes/services/userservice"
	"github.com/dgrijalva/jwt-go"
)

const oneWeek = time.Hour * 24 * 7
const credentialTimeoutDuration = oneWeek
const cerealNotesCookieName = "CerealNotesToken"

type JwtTokenClaim struct {
	models.UserId `json:"userId"`
	jwt.StandardClaims
}

var tokenSigningKey []byte

func SetTokenSigningKey(key []byte) {
	tokenSigningKey = key
}

// HANDLERS
func HandleLoginOrSignupRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
		// Check to see if user is already logged in, if so redirect
		if _, err := getUserIdFromJwtToken(request); err == nil {
			http.Redirect(
				responseWriter,
				request,
				paths.HomePath,
				http.StatusTemporaryRedirect)
			return
		}

		parsedTemplate, err := template.ParseFiles("templates/login_or_signup.tmpl")
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTemplate.Execute(responseWriter, nil)

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodGet})
	}
}

func HandleUserRequest(
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
		if err := userservice.StoreNewUser(
			signupForm.DisplayName,
			models.NewEmailAddress(signupForm.EmailAddress),
			signupForm.Password,
		); err != nil {
			if err == userservice.EmailAddressAlreadyInUseError {
				statusCode = http.StatusConflict
			} else {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			statusCode = http.StatusCreated
		}

		responseWriter.WriteHeader(statusCode)

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodPost})
	}
}

func RedirectRequestToHome(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
		http.Redirect(
			responseWriter,
			request,
			paths.HomePath,
			http.StatusTemporaryRedirect)
		return
	default:
		respondWithMethodNotAllowed(
			responseWriter,
			[]string{http.MethodGet})
	}
}

func HandleSessionRequest(
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

		if err := userservice.AuthenticateUserCredentials(
			models.NewEmailAddress(loginForm.EmailAddress),
			loginForm.Password,
		); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set our cookie to have a valid JWT Token as the value
		{
			userId, err := userservice.GetIdForUserWithEmailAddress(
				models.NewEmailAddress(loginForm.EmailAddress))
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			token, err := createTokenAsString(userId, credentialTimeoutDuration)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			expirationTime := time.Now().Add(credentialTimeoutDuration)

			cookie := http.Cookie{
				Name:     cerealNotesCookieName,
				Value:    token,
				Expires:  expirationTime,
				HttpOnly: true,
			}

			http.SetCookie(responseWriter, &cookie)
		}

		responseWriter.WriteHeader(http.StatusCreated)
		fmt.Fprint(responseWriter, "passward email combo was correct")

	case http.MethodDelete:
		// Cookie will overwrite existing cookie then delete itself
		cookie := http.Cookie{
			Name:     cerealNotesCookieName,
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		}

		http.SetCookie(responseWriter, &cookie)
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(fmt.Sprint("user succefully logged out")))
	default:
		respondWithMethodNotAllowed(
			responseWriter,
			[]string{http.MethodPost, http.MethodDelete})
	}
}

type AuthentictedRequestHandlerType func(
	http.ResponseWriter,
	*http.Request,
	models.UserId)

func AuthenticateOrRedirectToLogin(
	authenticatedHandlerFunc AuthentictedRequestHandlerType,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if userId, err := getUserIdFromJwtToken(request); err != nil {
			// If not loggedin redirect to login page
			http.Redirect(
				responseWriter,
				request,
				paths.LoginOrSignupPath,
				http.StatusTemporaryRedirect)
		} else {
			authenticatedHandlerFunc(responseWriter, request, userId)
		}
	}
}

// Authenticated Handlers
func HandleHomeRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles("templates/home.tmpl")
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		parsedTemplate.Execute(responseWriter, userId)
	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodGet})
	}
}

// UTIL
func respondWithMethodNotAllowed(
	responseWriter http.ResponseWriter,
	allowedMethods []string,
) {
	allowedMethodsString := strings.Join(allowedMethods, ", ")
	responseWriter.Header().Set("Allow", allowedMethodsString)

	statusCode := http.StatusMethodNotAllowed
	http.Error(responseWriter, http.StatusText(statusCode), statusCode)
}
