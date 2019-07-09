package DB

import (
	"database/sql"
	"log"
)

func ConnectDB(username, password, address, dbName string) *sql.DB {
	// Setup db connection
	db, err := sql.Open("mysql",
		username+":"+password+"@tcp("+address+")/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
