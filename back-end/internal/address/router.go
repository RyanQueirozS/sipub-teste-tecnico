package address

import (
	"net/http"
	"sipub-test/internal"
)

type UserRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewAddressRouter() UserRouter {
	router := UserRouter{
		controller: NewAddressController(),
	}
	return router
}

func (r UserRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/addresses"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r UserRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r UserRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r UserRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r UserRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r UserRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r UserRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("PUT "+r.baseEndPoint+"/{id}", r.controller.Update)
}
