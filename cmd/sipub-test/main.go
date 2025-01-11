package main

import (
	"net/http"
	internal "sipub-test/internal"
	product "sipub-test/internal/product"
)

const portNum string = ":8080"

// Golang doesn't allow slices to be constant, so var will do it
var routers = []internal.IRouter{
	product.ProductRouter{},
}

func main() {
	mux := http.NewServeMux()
	internal.RouterInitializeAll(mux, routers)
}
