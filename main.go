package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atmiguel/cerealnotes/handlers"
	"github.com/atmiguel/cerealnotes/models"
	"github.com/atmiguel/cerealnotes/routers"
)

// Get the current listening address
func determineListenPort() (string, error) {
	portEnvironmentVariableName := "PORT"
	port := os.Getenv(portEnvironmentVariableName)

	if len(port) == 0 {
		return "", fmt.Errorf(
			"environment variable %s not set",
			portEnvironmentVariableName)
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

func determineTokenSigningKey() ([]byte, error) {
	tokenSigningKeyVariableName := "TOKEN_SIGNING_KEY"
	tokenSigningKey := os.Getenv(tokenSigningKeyVariableName)

	if len(tokenSigningKey) == 0 {
		return nil, fmt.Errorf(
			"environment variable %s not set",
			tokenSigningKeyVariableName)
	}

	return []byte(tokenSigningKey), nil
}

func main() {
	// Set up db

	env := &handlers.Environment{}

	{
		databaseUrl, err := determineDatabaseUrl()
		if err != nil {
			log.Fatal(err)
		}

		db, err := models.ConnectToDatabase(databaseUrl, 20)
		if err != nil {
			log.Fatal(err)
		}

		env.Db = db

	}

	// Set up token signing key
	{
		tokenSigningKey, err := determineTokenSigningKey()
		if err != nil {
			log.Fatal(err)
		}
		env.TokenSigningKey = tokenSigningKey
	}

	// Start server
	{
		port, err := determineListenPort()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Listening on %s...\n", port)

		if err := http.ListenAndServe(port, routers.DefineRoutes(env)); err != nil {
			log.Fatal(err)
		}
	}
}
