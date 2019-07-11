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
	//getting shorten url from url
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
		var date string
		er := db.QueryRow("SELECT expire FROM urls WHERE shorten = '" + shortLink + "'").Scan(&date)
		if er != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
		if date < strings.Split(time.Now().String(), " ")[0] {
			w.Write([]byte("your link is expired."))
			log.Println(date)
		} else {
			//er := db.QueryRow("UPDATE urls SET cntr = cntr + 1 WHERE shorten = '" + shortLink +"'")
			//if er != nil && err != sql.ErrNoRows {
			//	log.Fatal(err)
			//}
			log.Println("request: " + shortLink + " => " + link)
			if link[0:1] != "//" && link[0:6] == "https:" && link[0:5] == "http:" {
				link = "//" + link
			}
			// do redirection
			http.Redirect(w, r, link, http.StatusTemporaryRedirect)
		}
	}
}

func GetURL(url string, short string, date string) (string, string) {
	//generate random shorten url
	if short == "" || len(short)-2 > len(url)/3 {
		fmt.Println("rand")
		var s string
		sh := 0
		for i := 0; i < len(url); i++ {
			sh += int(url[i])
		}
		fmt.Println(sh)
		ran := rand.New(rand.NewSource(time.Now().UnixNano()))
		short = "re" + strconv.Itoa(sh) + strconv.Itoa(ran.Intn(1000))
		e := db.QueryRow("SELECT shorten FROM urls WHERE shorten=" + short).Scan(&s)
		for !(e != nil && e == sql.ErrNoRows) {
			short = "re" + strconv.Itoa(sh) + strconv.Itoa(ran.Intn(1000))
			e = db.QueryRow("SELECT shorten FROM urls WHERE shorten=" + short).Scan(&s)
		}
	} else {
		var s string
		e := db.QueryRow("SELECT shorten FROM urls WHERE shorten=" + short).Scan(&s)
		if !(e != nil && e == sql.ErrNoRows) {
			fmt.Println("used : " + short)
			sh := 0
			for i := 0; i < len(url); i++ {
				sh += int(url[i])
			}
			fmt.Println(sh)
			ran := rand.New(rand.NewSource(time.Now().UnixNano()))
			short = "re" + strconv.Itoa(sh) + strconv.Itoa(ran.Intn(1000))
			e := db.QueryRow("SELECT shorten FROM urls WHERE shorten=" + short).Scan(&s)
			for !(e != nil && e == sql.ErrNoRows) {
				short = "re" + strconv.Itoa(sh) + strconv.Itoa(ran.Intn(1000))
				e = db.QueryRow("SELECT shorten FROM urls WHERE shorten=" + short).Scan(&s)
			}
		}
	}
	//default expire time
	if date == "" {
		dat := time.Now().AddDate(0, 0, expire)
		date = strings.Split(dat.String(), " ")[0]
	}
	log.Println("creat : " + url + " ==> " + short + " until : " + date)
	//save into database
	add, err := db.Prepare("INSERT INTO urls(shorten, url, expire) VALUES(\"" + short + "\",\"" + url + "\",'" + date + "')")
	if err != nil {
		log.Println(err.Error())
	}
	defer add.Close()
	if err != nil {
		log.Fatal(err)
	}
	add.Exec()
	return short, date
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
	router.HandleFunc(`/{url:re.*}`, Redirect).Methods("GET")
	router.HandleFunc(`/`, home)
	log.Fatal(http.ListenAndServe(":5000", router))
}

func home(w http.ResponseWriter, r *http.Request) {
	//var index template.HTML

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "server/index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		url := r.FormValue("url")
		short := "re" + r.FormValue("short")
		ex := r.FormValue("ex")
		var sh string
		sh, ex = GetURL(url, short, ex)
		if sh != short {
			if short == "" {
				fmt.Fprintf(w, "we generate a random shorten url for you\n")
			} else {
				if len(short)-2 > len(url)/3 {
					fmt.Fprintf(w, "your shorten url was too big ")
				} else {
					fmt.Fprintf(w, "your shorten url is already in use ")
				}
				fmt.Fprintf(w, "so we generate a random one\n")
			}
			short = sh
		}
		fmt.Fprintf(w, "%s/%s redirects you to %s until %s", r.Host, short, url, ex)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	//http.ServeFile(w, r, fmt.Sprintf("server/index.html"))
	//var url string
	//url = r.FormValue("uurl")
	//if url != "" {
	//	fmt.Println(url)
	//	sh := GetURL(w, url)
	//	r.Form.Set("short", sh)
	//
	//}
}
