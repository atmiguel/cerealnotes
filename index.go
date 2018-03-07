package main

import (
    "fmt"
    "log"
    "net/http"
)

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
    // set handler
    http.HandleFunc("/", rootHandler)

    // start server
    log.Fatal(http.ListenAndServe(":80", nil))
}
