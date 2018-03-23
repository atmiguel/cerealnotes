package main

import (
	"fmt"
	"github.com/atmiguel/cerealnotes/databaseutil"
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

func determineDatabaseUrl() (string, error) {
	environmentVariableName := "DATABASE_URL"
	databaseUrl := os.Getenv(environmentVariableName)

	if len(databaseUrl) == 0 {
		return "", fmt.Errorf(
			"environment variable %s not set",
			environmentVariableName)
	}

	return databaseUrl, nil
}

func main() {
	routers.SetRoutes()

	// SET UP DB
	{
		databaseUrl, err := determineDatabaseUrl()
		if err != nil {
			log.Fatal(err)
		}

		if err := databaseutil.ConnectToDatabase(databaseUrl); err != nil {
			log.Fatal(err)
		}
	}

	// START SERVER
	{
		port, err := determineListenPort()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Listening on %s...\n", port)

		if err := http.ListenAndServe(port, nil); err != nil {
			panic(err)
		}
	}
}
