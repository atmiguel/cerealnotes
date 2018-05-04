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
const baseTemplateName = "base"
const baseTemplateFile = "templates/base.tmpl"

// JwtTokenClaim contains all claims required for authentication, including the standard JWT claims.
type JwtTokenClaim struct {
	models.UserId `json:"userId"`
	jwt.StandardClaims
}

var tokenSigningKey []byte

func SetTokenSigningKey(key []byte) {
	tokenSigningKey = key
}

// UNAUTHENTICATED HANDLERS

// HandleLoginOrSignupRequest responds to unauthenticated GET requests with the login or signup page.
// For authenticated requests, it redirects to the home page.
func HandleLoginOrSignupRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
		if _, err := getUserIdFromJwtToken(request); err == nil {
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
		respondWithMethodNotAllowed(responseWriter, http.MethodPost)
	}
}

// HandleSessionRequest responds to POST requests by authenticating and responding with a JWT.
// It responds to DELETE requests by expiring the client's cookie.
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
			statusCode := http.StatusInternalServerError
			if err == userservice.CredentialsNotAuthorizedError {
				statusCode = http.StatusUnauthorized
			}
			http.Error(responseWriter, err.Error(), statusCode)
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
			http.MethodPost,
			http.MethodDelete)
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
			// If not logged in, redirect to login page
			http.Redirect(
				responseWriter,
				request,
				paths.LoginOrSignupPage,
				http.StatusTemporaryRedirect)
		} else {
			authenticatedHandlerFunc(responseWriter, request, userId)
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

func HandleHomeRequest(
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

func HandleNotesRequest(
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
	allowedMethods ...string,
) {
	if allowedMethods == nil {
		panic("This should never happen")
	}

	allowedMethodsString := strings.Join(allowedMethods, ", ")
	responseWriter.Header().Set("Allow", allowedMethodsString)

	statusCode := http.StatusMethodNotAllowed
	http.Error(responseWriter, http.StatusText(statusCode), statusCode)
}
