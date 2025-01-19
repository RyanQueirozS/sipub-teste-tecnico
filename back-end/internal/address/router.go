package address

import (
	"net/http"
	"sipub-test/internal"
)

type AddressRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewAddressRouter() AddressRouter {
	router := AddressRouter{
		controller: NewAddressController(),
	}
	return router
}

func (r AddressRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/addresses"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r AddressRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r AddressRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r AddressRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r AddressRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r AddressRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r AddressRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("PUT "+r.baseEndPoint+"/{id}", r.controller.Update)
}
