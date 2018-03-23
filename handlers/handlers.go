package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/atmiguel/cerealnotes/services/userservice"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

		var statusCode int
		if err := userservice.StoreNewUser(
			signupForm.DisplayName,
			signupForm.EmailAddress,
			signupForm.Password,
		); err != nil {
			if err == userservice.EmailAddressAlreadyInUseError {
				statusCode = http.StatusConflict
			} else {
				panic(err)
			}
		} else {
			statusCode = http.StatusCreated
		}

		responseWriter.WriteHeader(statusCode)

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
