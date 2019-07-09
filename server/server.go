package server

import (
	"fmt"
	"net/http"
	"strings"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	u := strings.Split(r.URL.String(), "/")
	fmt.Print(u[2])
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	u := strings.Split(r.URL.String(), "/")
	fmt.Print(u[2])
}
