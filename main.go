package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
)

func handleGetRootRequest(
        responseWriter http.ResponseWriter,
        request *http.Request) {

    request.ParseForm()

    fmt.Println(request.Form)
    fmt.Println("path", request.URL.Path)
    fmt.Println("scheme", request.URL.Scheme)
    fmt.Println(request.Form["url_long"])

    for key, value := range request.Form {
        fmt.Println("key:", key)
        fmt.Println("value:", strings.Join(value, ""))
    }

    fmt.Fprintf(responseWriter, "Hello World!")
}

func main() {
    http.HandleFunc("/", handleGetRootRequest) // set router

    // start server
    err := http.ListenAndServe(":8080", nil) // set listen port

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
