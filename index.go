package main

import (
    "fmt"
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

func rootHandler(
        responseWriter http.ResponseWriter,
        request *http.Request) {

    var responseContent string = fmt.Sprintf(
        "Hi there, I love %s!",
        request.URL.Path[1:])

    fmt.Fprintf(
        responseWriter,
        responseContent)
}

func main() {
    port, err := determineListenPort()
    if err != nil {
        log.Fatal(err)
    }

    // set handler
    http.HandleFunc("/", rootHandler)

    // start server
    log.Printf("Listening on %s...\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        panic(err)
    }
}
