package main

import (
	"encoding/json"
	"fmt"
	"github.com/atmiguel/cerealnotes/services/userservice"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// get the current listening address or fail if input is not correct
func determineListenPort() (string, error) {
	port := os.Getenv("PORT")

	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}

	return ":" + port, nil
}

func respondWithMethodNotAllowed(
	responseWriter http.ResponseWriter,
	allowedMethods []string,
) {
	allowedMethodsString := strings.Join(allowedMethods, ", ")
	responseWriter.Header().Set("Allow", allowedMethodsString)

	statusCode := http.StatusMethodNotAllowed
	http.Error(responseWriter, http.StatusText(statusCode), statusCode)
}

func handleLoginOrSignupRequest(
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

func handleUserRequest(
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

		_, err := userservice.CreateUser(
			signupForm.DisplayName,
			signupForm.EmailAddress,
			signupForm.Password)
		if err != nil {
			panic(err)
		}

		// TODO js should check status returned
		responseWriter.WriteHeader(http.StatusCreated)

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodPost})
	}
}

func handleSessionRequest(
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

		if err := userservice.AuthenticateUser(
			loginForm.EmailAddress,
			loginForm.Password,
		); err != nil {
			// TODO check err type
			panic(err)
		}

		responseWriter.WriteHeader(http.StatusCreated)
		responseWriter.Write([]byte(fmt.Sprint("passward email combo was correct")))

	default:
		respondWithMethodNotAllowed(responseWriter, []string{http.MethodPost})
	}
}

func main() {
	// SET ROUTER

	// static files
	{
		staticDirectoryName := "static"
		staticDirectoryPaddedWithSlashes := "/" + staticDirectoryName + "/"

		fileServer := http.FileServer(http.Dir(staticDirectoryName))

		http.Handle(
			staticDirectoryPaddedWithSlashes,
			http.StripPrefix(staticDirectoryPaddedWithSlashes, fileServer))
	}

	// templates
	http.HandleFunc("/login-or-signup", handleLoginOrSignupRequest)

	// forms
	http.HandleFunc("/user", handleUserRequest)
	http.HandleFunc("/session", handleSessionRequest)

	// START SERVER
	port, err := determineListenPort()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
