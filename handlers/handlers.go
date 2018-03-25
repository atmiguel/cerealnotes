package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/services/userservice"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type CerealNotesClaims struct {
	UserId models.UserId `json:"UserId"`
	jwt.StandardClaims
}

var tokenSigningKey []byte

func SetSigningKey(key []byte) {
	tokenSigningKey = key
}

// HANDLERS
func HandleLoginOrSignupRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
		// Check to see if they are already logged in if so redirect 
		if _, err := getUserIdFromStoredToken(request); err == nil {
			http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
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

		// Create token as a cookie and set it
		{
			userId, err := userservice.GetUserIdFromEmailAddress(loginForm.EmailAddress)
			if err != nil {
				panic(err)
			}

			token, err := createTokenAsString(userId, oneWeekInMinutes)
			if err != nil {
				panic(err)
			}

			expiration := time.Now().Add(oneWeekInMinutes * time.Minute)

			cookie := http.Cookie{
				Name:     cerealNotesCookieName,
				Value:    token,
				Expires:  expiration,
				HttpOnly: true,
			}

			http.SetCookie(responseWriter, &cookie)

		}

		responseWriter.WriteHeader(http.StatusCreated)
		responseWriter.Write([]byte(fmt.Sprint("passward email combo was correct")))

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodPost})
	}
}

type AuthentictedRequestHandlerType func(http.ResponseWriter, *http.Request, models.UserId)


func AuthenticateOrRedirectToLogin(
	originalHandlerFunc AuthentictedRequestHandlerType,
) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if userId, err := getUserIdFromStoredToken(request); err == nil {
			originalHandlerFunc(responseWriter, request, userId) // call original
		} else {
			http.Redirect(
				responseWriter, 
				request, 
				"/login-or-signup", 
				http.StatusSeeOther,
			)
		}
	})
}

func HandleRootRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
	userId models.UserId,
) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles("templates/root.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		parsedTemplate.Execute(responseWriter, nil)
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

func parseTokenFromString(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(
		strings.TrimSpace(tokenString),
		&CerealNotesClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenSigningKey, nil
		})
	return token, err
}

// one week
const oneWeekInMinutes = 60 * 24 * 7
const cerealNotesCookieName = "CerealNotesToken"

func createTokenAsString(
	userId models.UserId,
	expirationTimeInMinutes int64,
) (string, error) {
	// Create the Claims
	claims := CerealNotesClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + (expirationTimeInMinutes * 60),
			Issuer:    "CerealNotes",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(tokenSigningKey)
	return ss, err
}

func getUserIdFromStoredToken(request *http.Request) (models.UserId, error) {
	cookie, err := request.Cookie(cerealNotesCookieName)
	if err != nil {
		return 0, err
	}

	token, err := parseTokenFromString(cookie.Value)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*CerealNotesClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return 0, errors.Errorf("Token was invalid or unreadable")
	}
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
	if claims, ok := token.Claims.(*CerealNotesClaims); ok && token.Valid {
		if claims.UserId != 32 {
			log.Fatal("error in token")
		}
		fmt.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println("Token claims could not be read")
	}
}
