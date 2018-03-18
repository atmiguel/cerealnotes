package main

import (
	"fmt"
	"github.com/atmiguel/cerealnotes/routers"
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

func main() {
	routers.SetRoutes()

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
