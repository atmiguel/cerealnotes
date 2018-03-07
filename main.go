package main

import (
    "html/template"
    "log"
    "net/http"
)

func handleGetRootRequest(
        responseWriter http.ResponseWriter,
        request *http.Request) {

    // TODO handle error
    parsedTemplate, _ := template.ParseFiles("root.tmpl")
    parsedTemplate.Execute(responseWriter, nil)
}

func main() {
    http.HandleFunc("/", handleGetRootRequest) // set router

    // start server
    err := http.ListenAndServe(":8080", nil) // set listen port

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
