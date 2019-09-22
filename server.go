package main

import (
	"net/http"
	"fmt"
)

func main() {
	fmt.Println("Server started")
	http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./public`)))
}