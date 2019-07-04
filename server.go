package main

import (
	"net/http"
)

func HandleRequests() {
	server := http.Server{
		Addr: "urlshortner:8090",
	}
	http.HandleFunc("/", homePage)
	server.ListenAndServe()
}

func homePage(w http.ResponseWriter, r *http.Request) {

}
