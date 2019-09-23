package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server starting")
	http.ListenAndServe(`:8080`, http.Handler(http.FileServer(http.Dir("./public"))))
}
