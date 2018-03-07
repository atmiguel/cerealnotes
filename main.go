package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "strings"
)

func determineListenPort() (string, error) {
  port := os.Getenv("PORT")

  if port == "" {
    return "", fmt.Errorf("$PORT not set")
  }

  return ":" + port, nil
}

func respondWithMethodNotAllowed(
        responseWriter http.ResponseWriter,
        allowedMethods []string) {

    statusCode := http.StatusMethodNotAllowed

    responseWriter.Header().Set(
        "Allow",
        strings.Join(allowedMethods, ", "))

    http.Error(
        responseWriter,
        http.StatusText(statusCode),
        statusCode)
}

func handleLoginOrSignupRequest(
        responseWriter http.ResponseWriter,
        request *http.Request) {

    switch request.Method {
        case http.MethodGet:
            parsedTemplate, err := template.ParseFiles("templates/login_or_signup.tmpl")
            if err != nil {
                log.Fatal(err) // TODO is this right?
            }

            parsedTemplate.Execute(responseWriter, nil)

        default:
            respondWithMethodNotAllowed(
                responseWriter,
                []string{http.MethodGet})
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

    // templates
    http.HandleFunc("/login-or-signup", handleLoginOrSignupRequest)

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
