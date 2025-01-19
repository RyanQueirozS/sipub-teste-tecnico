package shopping_cart

import (
	"net/http"
	"sipub-test/internal"
)

type ShoppingCartRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewShoppingCartRouter() ShoppingCartRouter {
	router := ShoppingCartRouter{
		controller: NewShoppingCartController(),
	}
	return router
}

func (r ShoppingCartRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/cart"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r ShoppingCartRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r ShoppingCartRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r ShoppingCartRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r ShoppingCartRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r ShoppingCartRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r ShoppingCartRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("PUT "+r.baseEndPoint+"/{id}", r.controller.Update)
}
