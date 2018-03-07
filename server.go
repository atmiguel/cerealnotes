package main

import (
    "html/template"
    "log"
    "net/http"
)

func rootHandler(
        responseWriter http.ResponseWriter,
        request *http.Request) {

    // TODO handle error
    parsedTemplate, _ := template.ParseFiles("root.html")
    parsedTemplate.Execute(responseWriter, nil)
}

func main() {
    // set handler
    http.HandleFunc("/", rootHandler)

    // start server
    log.Fatal(http.ListenAndServe(":8080", nil))
}
