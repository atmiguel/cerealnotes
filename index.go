package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
)

// get the current listening address or fail if input is not correct
func determineListenAddress() (string, error) {
  portStr := os.Getenv("PORT")
  
  if portStr == "" {
    return "", fmt.Errorf("$PORT not set")
  }
  
  port, err := strconv.Atoi(portStr)
  if err == nil {
    return "", fmt.Errorf("$PORT not a valid integer")
  }
  return fmt.Sprintf(":%d", port), nil
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
    err := http.ListenAndServe(addr, nil)
    if err != nil {
        panic(err)
    }
}
