package server

import "fmt"

func MakeHash(url string) {
	h := 0
	for i := 0; i < len(url); i++ {
		h += int(url[i])
	}
	fmt.Print(h)
}
