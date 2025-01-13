package main

import (
	"net/http"
	internal "sipub-test/internal"
	product "sipub-test/internal/product"
)

const portNum string = ":8080"

// Loops through an array of routers and uses the Init method on them.
// The `Init` method should initialize all of the methods per route.
func RouterInitializeAll(mux *http.ServeMux, routers []internal.IRouter) {
	for i := 0; i < len(routers); i++ {
		routers[i].Init(mux)
	}
}

// Golang doesn't allow slices to be constant, so var will do it
var routers = []internal.IRouter{
	product.ProductRouter{},
}

func main() {
	mux := http.NewServeMux()
	RouterInitializeAll(mux, routers)
}
