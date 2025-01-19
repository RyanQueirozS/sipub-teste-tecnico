package main

import (
	"log"
	"net/http"
	"sipub-test/db"
	internal "sipub-test/internal"
	"sipub-test/internal/address"
	"sipub-test/internal/delivery"
	"sipub-test/internal/delivery_product"
	"sipub-test/internal/payment"
	"sipub-test/internal/product"
	"sipub-test/internal/shopping_cart"
	"sipub-test/internal/user"
	"sipub-test/internal/user_address"
	"sipub-test/internal/user_delivery"

	"github.com/rs/cors"
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

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})
	mux := http.NewServeMux()
	RouterInitializeAll(mux,
		address.NewAddressRouter(),
		delivery.NewDeliveryRouter(),
		delivery_product.NewDeliveryProductRouter(),
		payment.NewPaymentRouter(),
		product.NewProductRouter(),
		shopping_cart.NewShoppingCartRouter(),
		user.NewUserRouter(),
		user_address.NewUserAddressRouter(),
		user_delivery.NewUserDeliveryRouter(),
	)
	handler := corsHandler.Handler(mux)

	log.Println("Starting server...")
	http.ListenAndServe(portNum, handler)
}
