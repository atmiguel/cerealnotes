package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func determineListenAddress() (string, error) {
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
    addr, err := determineListenAddress()
    if err != nil {
        log.Fatal(err)
    }

    // set handler
    http.HandleFunc("/", rootHandler)

    // start server
    log.Printf("Listening on %s...\n", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        panic(err)
    }
}
