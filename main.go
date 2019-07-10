package main

import (
	"../URLshortener/DB"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var db sql.DB

const expire = 20

func Redirect(w http.ResponseWriter, r *http.Request) {
	//u := strings.Split(r.URL.String(), "/")
	//fmt.Println(u[1])
	params := mux.Vars(r)
	shortLink := params["url"]
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
		if link[0:1] != "//" && link[0:6] == "https:" && link[0:5] == "http:" {
			link = "//" + link
		}
		// do redirection
		http.Redirect(w, r, link, http.StatusTemporaryRedirect)
	}
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	u := strings.Split(r.URL.String(), "/")
	fmt.Println(u[2])
	sh := 0
	for i := 0; i < len(u[2]); i++ {
		sh += int(u[2][i])
	}
	fmt.Println(sh)
	short := "re" + strconv.Itoa(sh) + strconv.Itoa(rand.Intn(100))
	dat := time.Now().AddDate(0, 0, expire)
	date := strings.Split(dat.String(), " ")[0]
	log.Println(u[2] + " ==> " + short + " until : " + date)
	add, err := db.Prepare("INSERT INTO urls(shorten, url, expire) VALUES(\"" + short + "\",\"" + u[2] + "\",'" + date + "')")
	if err != nil {
		log.Println(err.Error())
	}
	defer add.Close()
	if err != nil {
		log.Fatal(err)
	}
	add.Exec()

	w.Write([]byte("yor shorten url : " + short))
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
	log.Println("Server available")
	router.HandleFunc(`/get/{url:.*}`, GetURL).Methods("GET")
	router.HandleFunc(`/{url:re.*}`, Redirect).Methods("GET")
	router.HandleFunc(`/`, home)
	log.Fatal(http.ListenAndServe(":5000", router))
}
func home(w http.ResponseWriter, r *http.Request) {
	//var index template.HTML
	http.ServeFile(w, r, fmt.Sprintf("server/index.html"))
	url := r.FormValue("uurl")
	if url != "" {
		fmt.Println(url)
		http.Redirect(w, r, "get/"+url, http.StatusSeeOther)
	}
}
