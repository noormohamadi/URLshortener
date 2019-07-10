package main

import (
	"../URLshortener/DB"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

var db sql.DB

func Redirect(w http.ResponseWriter, r *http.Request) {
	//u := strings.Split(r.URL.String(), "/")
	//fmt.Println(u[1])
	params := mux.Vars(r)
	shortLink := params["url"]
	fmt.Println(shortLink)
	var link string
	err := db.QueryRow("SELECT url FROM urls WHERE shorten ='" + shortLink + "'").Scan(&link)

	// Check for errors or invalid short urls
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	} else if err != nil && err == sql.ErrNoRows {
		log.Println(shortLink + " not found.")
		w.Write([]byte("404 Not found."))
	} else {
		// Valid short url
		log.Println("request: " + shortLink + " => " + link)
		// do redirection
		http.Redirect(w, r, link, http.StatusTemporaryRedirect)
	}
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	u := strings.Split(r.URL.String(), "/")
	fmt.Println(u[2])
}

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
	router.HandleFunc(`/get/{url:.*}`, GetURL).Methods("GET")
	router.HandleFunc(`/{url:re.*}`, Redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}
