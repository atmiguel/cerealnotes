package main

import (
	"CerealNotes/databaseutil"
	"encoding/json"
	"fmt"
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

func respondWithMethodNotAllowed(responseWriter http.ResponseWriter, allowedMethods []string) {
	statusCode := http.StatusMethodNotAllowed

	responseWriter.Header().Set(
		"Allow",
		strings.Join(allowedMethods, ", "))

	http.Error(
		responseWriter,
		http.StatusText(statusCode),
		statusCode)
}

func handleLoginOrSignupRequest(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		parsedTemplate, err := template.ParseFiles("templates/login_or_signup.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		parsedTemplate.Execute(responseWriter, nil)

	default:
		respondWithMethodNotAllowed(
			responseWriter,
			[]string{http.MethodGet})
	}
}

// TODO cleanup error cases
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

type UserId struct {
	Value int64 `json:"value"`
}

func handleUserRequest(responseWriter http.ResponseWriter, request *http.Request) {
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

		val, err := databaseutil.SaveNewUser(
			signupForm.DisplayName,
			signupForm.EmailAddress,
			signupForm.Password)
		if err != nil {
			panic(err)
		}

		// TODO create User
		userId := UserId{Value: val}

		responseWriter.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(responseWriter).Encode(userId); err != nil {
			panic(err)
		}

	default:
		respondWithMethodNotAllowed(
			responseWriter,
			[]string{http.MethodPost})
	}
}

func handleSessionRequest(
	responseWriter http.ResponseWriter,
	request *http.Request) {

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

		validated, err := databaseutil.ValidateUser(
			loginForm.EmailAddress,
			loginForm.Password)

		if err != nil {
			panic(err)
		}

		log.Printf("did we find the user + password combo in the table: %t", validated)
		responseWriter.WriteHeader(http.StatusCreated)
		responseWriter.Write([]byte(fmt.Sprintf("passward email combo was correct? %t", validated)))

	default:
		respondWithMethodNotAllowed(
			responseWriter,
			[]string{http.MethodPost})
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
			http.StripPrefix(
				staticDirectoryPaddedWithSlashes,
				fileServer))
	}

	err := databaseutil.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
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
