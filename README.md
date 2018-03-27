# Installation
## Locally
* postgres server installed and running
	* `brew install postgres`
	* `pg_ctl -D /usr/local/var/postgres start`
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
2. `heroku run local`
3. Visit `localhost:8080/login-or-signup`

# Run DB migrations
More db information in `migrations/README.md`