package databaseutil

import (
	"database/sql"
	// Notice that we’re loading the driver anonymously, The driver registers itself as being available to the database/sql package.
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
	"fmt"
)

type Dbutil struct {
	db *sql.DB
}

func (wrapper Dbutil) Connect(dbUrl string) (error) {
	temp, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}

	wrapper.db = temp

	// Quickly test if the connection to the database worked.
	if err := wrapper.db.Ping(); err != nil {
		return err
	}

	return nil
}


func (wrapper Dbutil) SaveNewUser(
	displayName string, 
	emailAddress string, 
	password string) (int64, error) {
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	sqlStatement := `
		INSERT INTO users (display_name, email_address, password, creation_time)
		VALUES ($1, $2, $3, $4)`

	result, err := wrapper.db.Exec(sqlStatement, displayName, emailAddress, hash, time.Now().UTC())
	if err != nil {
		return -1, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, nil
	}

	fmt.Sprintf("%d", id)
	return id, nil
}

func (wrapper Dbutil) validateUser(
	emailAddress string,
	password string) (bool, error){

	// todo get from database
	hashFromDatabase := []byte(emailAddress)

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}
