package server

import (
	"fmt"
	"html/template"
	"net/http"
)

var index template.HTML

func HandleRequests() {
	server := http.Server{
		Addr: "urlshortner:8090",
	}
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/*", Redirect)

	server.ListenAndServe()
}

func HomePage(w http.ResponseWriter, r *http.Request) {

}
func Redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Print("hiiiiiiiiiiiiiiiiiiii")
}
