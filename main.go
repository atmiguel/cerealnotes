package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
)

// get the current listening address or fail if input is not correct
func determineListenPort() (string, error) {
    port := os.Getenv("PORT")

    if port == "" {
        return "", fmt.Errorf("$PORT not set")
    }

    return ":" + port, nil
}

func handleGetRootRequest(
        responseWriter http.ResponseWriter,
        request *http.Request) {

    // TODO handle error
    parsedTemplate, _ := template.ParseFiles("root.tmpl")
    parsedTemplate.Execute(responseWriter, nil)
}

func main() {
    port, err := determineListenPort()
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/", handleGetRootRequest) // set router

    // start server
    log.Printf("Listening on %s...\n", port)

    if err := http.ListenAndServe(port, nil); err != nil {
        panic(err)
    }
}
