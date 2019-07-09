package main

import (
	"../URLshortener/DB"
	"../URLshortener/server"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var db sql.DB

func main() {
	router := mux.NewRouter()
	db = *DB.ConnectDB("root", "", "localhost:3306", "urls")
	defer db.Close()
	err := db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	router.SkipClean(true) // keep double slash for URLs like http://...
	log.Println("Server starting...")
	router.HandleFunc(`/get/{url:.*}`, server.GetURL).Methods("GET")
	router.HandleFunc(`/re/{url:.*}`, server.Redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}
