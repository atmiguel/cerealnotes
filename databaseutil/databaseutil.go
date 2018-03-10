package databaseutil

import (
	"database/sql"
	// Notice that we’re loading the driver anonymously, The driver registers itself as being available to the database/sql package.
	_ "github.com/lib/pq"
)

func Connect(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	// Quickly test if the connection to the database worked.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
