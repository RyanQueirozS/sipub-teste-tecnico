package user_delivery

import (
	"net/http"
	"sipub-test/internal"
)

type UserDeliveryRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewUserDeliveryRouter() UserDeliveryRouter {
	router := UserDeliveryRouter{
		controller: NewUserDeliveryController(),
	}
	return router
}

func (r UserDeliveryRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/user_delivery"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r UserDeliveryRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r UserDeliveryRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r UserDeliveryRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r UserDeliveryRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r UserDeliveryRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r UserDeliveryRouter) update(mux *http.ServeMux) {
	// This needs to be implemented because of the interface, altough it won't
	// be used since the delivery-product shouldn't be updated
}
