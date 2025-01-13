package main

import (
	"log"
	"net/http"
	internal "sipub-test/internal"
	product "sipub-test/internal/product"
)

const portNum string = ":8080"

// Loops through an array of routers and uses the Init method on them.
// The `Init` method should initialize all of the methods per route.
func RouterInitializeAll(mux *http.ServeMux, routers ...internal.IRouter) {
	for i := 0; i < len(routers); i++ {
		routers[i].Init(mux)
	}
}

func main() {
	mux := http.NewServeMux()
	RouterInitializeAll(mux, product.NewProductRouter())
	log.Println("Starting server...")
	http.ListenAndServe(portNum, mux)
}
