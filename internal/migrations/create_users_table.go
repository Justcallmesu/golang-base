package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateUsersTable(database *sql.DB, channel chan bool) {
	const Create_User_Table = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER AUTO_INCREMENT,
		email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
		PRIMARY KEY(id)
	)
	`
	_, initializeError := database.Exec(Create_User_Table)

	if initializeError != nil {
		log.Fatal(initializeError)
	}

	fmt.Println("Table Initializer: User Table Initialized")
	channel <- true
	close(channel)

}
