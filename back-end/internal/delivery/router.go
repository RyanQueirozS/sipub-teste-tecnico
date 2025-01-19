package delivery

import (
	"net/http"
	"sipub-test/internal"
)

type DeliveryRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewDeliveryRouter() DeliveryRouter {
	router := DeliveryRouter{
		controller: NewDeliveryController(),
	}
	return router
}

func (r DeliveryRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/deliveries"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r DeliveryRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r DeliveryRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r DeliveryRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r DeliveryRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r DeliveryRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r DeliveryRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("PUT "+r.baseEndPoint+"/{id}", r.controller.Update)
}
