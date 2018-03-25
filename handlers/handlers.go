package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/atmiguel/cerealnotes/services/userservice"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Todo this should be pulled from an environemnt variable or something
var tokenSigningKey string = "bob"

// HANDLERS
func HandleLoginOrSignupRequest(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	switch request.Method {
	case http.MethodGet:
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

		responseWriter.WriteHeader(http.StatusCreated)
		responseWriter.Write([]byte(fmt.Sprint("passward email combo was correct")))

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodPost})
	}
}

type CerealNotesClaims struct {
	UserId string `json:"UserId"`
	jwt.StandardClaims
}

func createClaim(userId string) (string, error) {
	// Create the Claims
	claims := CerealNotesClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(tokenSigningKey)
	return ss, err
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

func parseTokenString(tokenString string) (*jwt.Token, error){
	// Parse the token
	token, err := jwt.ParseWithClaims(
		strings.TrimSpace(tokenString),
		&CerealNotesClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return tokenSigningKey, nil
		})
	return token, err
}
