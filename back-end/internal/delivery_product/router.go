package delivery_product

import (
	"net/http"
	"sipub-test/internal"
)

type DeliveryProductRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewDeliveryProductRouter() DeliveryProductRouter {
	router := DeliveryProductRouter{
		controller: NewDeliveryProductController(),
	}
	return router
}

func (r DeliveryProductRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/delivery_product"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r DeliveryProductRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r DeliveryProductRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r DeliveryProductRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r DeliveryProductRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r DeliveryProductRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r DeliveryProductRouter) update(mux *http.ServeMux) {
	// This needs to be implemented because of the interface, altough it won't
	// be used since the delivery-product shouldn't be updated
}
