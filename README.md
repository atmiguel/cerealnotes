# Installation
## Locally
* postgres server installed and running
    * `brew install postgres`
    * `pg_ctl start -D /usr/local/var/postgres`
* golang installed
    * `brew install go`
* godep installed:
    * `go get github.com/tools/godep`
* create database
    * `psql -d postgres -c "create database cerealnotes_test"`
    * Refer to `migrations/README.md` to create relevant database tables
* export environment variables
    ```
    export DATABASE_URL=postgres://localhost:5432/cerealnotes_test?sslmode=disable
    export TOKEN_SIGNING_KEY=abcdefg
    export PORT=8000
    ```
* run the app
    * `go build && ./cerealnotes`

## Heroku
* heroku cli installed
    * `brew install heroku`
* heroku instance

# Running CerealNotes

Assuming your local environment is setup correctly with Golang standards, you can start your local server with the following commands

1. `cd to this repo`
2. `go install && heroku local`
3. Visit `localhost:8080/login-or-signup`

# Run DB migrations
More db information in `migrations/README.md`
