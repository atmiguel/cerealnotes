package cerealNotesDb

import (
	"database/sql"
	// Notice that we’re loading the driver anonymously, The driver registers itself as being available to the database/sql package.
	_ "github.com/lib/pq"
	"log"
)

func Connect(dbUrl string) (*sql.DB) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

	return db
}
