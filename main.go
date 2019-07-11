package main

import (
	"../URLshortener/DB"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var db sql.DB

const expire = 20

type htmlvalues struct {
	Massage string
	Shorten string
	URL     string
	Expire  string
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	//getting shorten url from url
	params := mux.Vars(r)
	shorten := params["url"]
	if !DB.Exist(shorten, &db) {
		log.Println(shorten + " not found.")
		w.Write([]byte("404 Not found."))
	} else {
		// Valid short url
		link := DB.Select("url", shorten, &db)
		date := DB.Select("expire", shorten, &db)
		if date < strings.Split(time.Now().String(), " ")[0] {
			w.Write([]byte("your link is expired."))
			log.Println(date)
		} else {
			log.Println("request: " + shorten + " => " + link)
			if link[0:1] != "//" && link[0:6] == "https:" && link[0:5] == "http:" {
				link = "//" + link
			}
			DB.Used(shorten, &db)
			// do redirection
			http.Redirect(w, r, link, http.StatusTemporaryRedirect)
		}
	}
}

func GetURL(url string, short string, date string) (string, string) {
	//generate random shorten url
	if short == "" {
		log.Println("generate random shorten url ...")
		sh := 0
		for i := 0; i < len(url); i++ {
			sh += int(url[i])
		}
		size := 0
		for i := sh; i > 0; i /= 10 {
			size++
		}
		ran := rand.New(rand.NewSource(time.Now().UnixNano()))
		sh -= rand.Intn(sh)
		size = len(url)/3 - size
		if size <= 0 {
			sh /= int(math.Pow(10, math.Abs(float64(size))+1))
			size = 1
		}
		short = "re" + strconv.Itoa(sh) + strconv.Itoa(ran.Intn(int(math.Pow(10, float64(size)))))
		for DB.Exist(short, &db) {
			fmt.Println("hi")
			short = "re" + strconv.Itoa(sh) + strconv.Itoa(ran.Intn(int(math.Pow(10, float64(size)))))
		}
	}
	//default expire time
	if date == "" {
		dat := time.Now().AddDate(0, 0, expire)
		date = strings.Split(dat.String(), " ")[0]
	}
	log.Println("creat : " + url + " ==> " + short + " until : " + date)
	//save into database
	DB.Add(short, url, date, &db)
	return short, date
}

func home(w http.ResponseWriter, r *http.Request) {
	//var index template.HTML

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "server/index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		massage := ""
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		url := r.FormValue("url")
		short := "re" + r.FormValue("short")
		ex := r.FormValue("ex")
		if short == "re" {
			massage += "we generate a random shorten url for you !"
			short = ""
		} else {
			if len(short)-2 > len(url)/3 {
				size := len(url) / 3
				massage += "your shorten url was too big (it must be at last " + strconv.Itoa(size) + " characters) so we generate a random one !"
				short = ""
			} else {

				if DB.Exist(short, &db) {
					massage += short + " is already in use so we generate a random one !"
					short = ""
				}
			}
		}
		short, ex = GetURL(url, short, ex)
		short = r.Host + "/" + short
		value := htmlvalues{massage, short, url, ex}
		fp := path.Join("server", "shorten.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

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
