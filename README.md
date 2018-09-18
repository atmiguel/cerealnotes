# Installation
## Locally
* postgres server installed and running: please refer to `migrations/README.md` for more info
* heroku cli installed
	* `brew install heroku`
* golang installed
	* `brew install go`
* godep installed: 
	* `go get github.com/tools/godep`

## Heroku
* heroku instance
* instance connected to postgres db

# Running CerealNotes

Assuming your local environment is setup correctly with Golang standards, you can start your local server with the following commands

1. `cd to this repo`
2. `go install && heroku local`
3. Visit `localhost:8080/`
