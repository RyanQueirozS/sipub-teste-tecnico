package main

import (
	"fmt"
	"net/http"
)

const portNum string = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page")
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", Home)
}
