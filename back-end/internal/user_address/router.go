package user_address

import (
	"net/http"
	"sipub-test/internal"
)

type UserAddressRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewUserAddressRouter() UserAddressRouter {
	router := UserAddressRouter{
		controller: NewUserAddressController(),
	}
	return router
}

func (r UserAddressRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/user_address"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
}

func (r UserAddressRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r UserAddressRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r UserAddressRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r UserAddressRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r UserAddressRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}
