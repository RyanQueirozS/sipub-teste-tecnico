package main

import (
	"log"
	"net/http"
	"sipub-test/db"
	internal "sipub-test/internal"
	"sipub-test/internal/address"
	product "sipub-test/internal/product"
	"sipub-test/internal/user"
)

const portNum string = ":8080"

// Loops through an 'array' of routers and uses the Init method on them.
// The `Init` method should initialize all of the methods per route.
func RouterInitializeAll(mux *http.ServeMux, routers ...internal.IRouter) {
	for i := 0; i < len(routers); i++ {
		routers[i].Init(mux)
	}
}

func main() {
	dsn := "user:password@tcp(mysql_db:3306)/sipub_test"
	if err := db.InitializeDB(dsn); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	mux := http.NewServeMux()
	RouterInitializeAll(mux,
		product.NewProductRouter(),
		user.NewUserRouter(),
		address.NewAddressRouter(),
	)
	log.Println("Starting server...")
	http.ListenAndServe(portNum, mux)
}
