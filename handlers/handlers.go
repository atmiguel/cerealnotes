package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/paths"
	"github.com/atmiguel/cerealnotes/services/userservice"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type JwtTokenClaim struct {
	userId models.UserId `json:"UserId"`
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
			http.Redirect(responseWriter, request, paths.HomePath, http.StatusTemporaryRedirect)
			return
		}

		parsedTemplate, err := template.ParseFiles("templates/login_or_signup.tmpl")
		if err != nil {
			log.Fatal(err)
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
		var signupForm SignupForm
		body := getRequestBody(request)

		if err := json.Unmarshal(body, &signupForm); err != nil {
			panic(err)
		}

		err := userservice.StoreNewUser(
			signupForm.DisplayName,
			signupForm.EmailAddress,
			signupForm.Password)
		if err != nil {
			panic(err)
		}

		responseWriter.WriteHeader(http.StatusCreated)

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodPost})
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
		var loginForm LoginForm
		body := getRequestBody(request)

		if err := json.Unmarshal(body, &loginForm); err != nil {
			panic(err)
		}

		if err := userservice.AuthenticateUserCredentials(
			loginForm.EmailAddress,
			loginForm.Password,
		); err != nil {
			panic(err)
		}

		// Set cerealNotesCookieName to have value as a our JWT Token
		{
			userId, err := userservice.GetIdForUserWithEmailAddress(loginForm.EmailAddress)
			if err != nil {
				panic(err)
			}

			token, err := createTokenAsString(userId, oneWeekInMinutes)
			if err != nil {
				panic(err)
			}

			expirationTime := time.Now().Add(oneWeekInMinutes * time.Minute)

			cookie := http.Cookie{
				Name:     cerealNotesCookieName,
				Value:    token,
				Expires:  expirationTime,
				HttpOnly: true,
			}

			http.SetCookie(responseWriter, &cookie)
		}

		responseWriter.WriteHeader(http.StatusCreated)
		responseWriter.Write([]byte(fmt.Sprint("passward email combo was correct")))

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
) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			if userId, err := getUserIdFromJwtToken(request); err != nil {
				// If not loggedin redirect to login page
				http.Redirect(
					responseWriter,
					request,
					paths.LoginOrSignupPath,
					http.StatusSeeOther)
			} else {
				authenticatedHandlerFunc(responseWriter, request, userId)
			}
		})
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
			log.Fatal(err)
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

func getRequestBody(request *http.Request) []byte {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	if err := request.Body.Close(); err != nil {
		panic(err)
	}

	return body
}

// Token Util
const oneWeekInMinutes = 60 * 24 * 7
const cerealNotesCookieName = "CerealNotesToken"

func parseTokenFromString(tokenAsString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		strings.TrimSpace(tokenAsString),
		&JwtTokenClaim{},
		func(*jwt.Token) (interface{}, error) {
			return tokenSigningKey, nil
		})
}

func createTokenAsString(
	userId models.UserId,
	expirationTimeInMinutes int64,
) (string, error) {
	claims := JwtTokenClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + (expirationTimeInMinutes * 60),
			Issuer:    "CerealNotes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenSigningKey)
}

func getUserIdFromJwtToken(request *http.Request) (models.UserId, error) {
	cookie, err := request.Cookie(cerealNotesCookieName)
	if err != nil {
		return -1, err
	}

	token, err := parseTokenFromString(cookie.Value)
	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(*JwtTokenClaim); ok && token.Valid {
		return claims.userId, nil
	}
	return -1, errors.Errorf("Token was invalid or unreadable")
}

func tokenTest1() {
	var num models.UserId = 32
	bob, err := createTokenAsString(num, 1)
	if err != nil {
		fmt.Println("create error")
		log.Fatal(err)
	}

	token, err := parseTokenFromString(bob)
	if err != nil {
		fmt.Println("parse error")
		log.Fatal(err)
	}
	fmt.Println(bob)
	if claims, ok := token.Claims.(*JwtTokenClaim); ok && token.Valid {
		if claims.userId != 32 {
			log.Fatal("error in token")
		}
		fmt.Printf("%v %v", claims.userId, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println("Token claims could not be read")
	}
}
