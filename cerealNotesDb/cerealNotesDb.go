package cerealNotesDb

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func Connect(dbUrl string) (*sql.DB) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
