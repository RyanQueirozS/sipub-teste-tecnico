package payment

import (
	"net/http"
	"sipub-test/internal"
)

type PaymentRouter struct {
	baseEndPoint string
	controller   internal.IController
}

func NewPaymentRouter() PaymentRouter {
	router := PaymentRouter{
		controller: NewPaymentController(),
	}
	return router
}

func (r PaymentRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/payment"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r PaymentRouter) create(mux *http.ServeMux) {
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r PaymentRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r PaymentRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r PaymentRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r PaymentRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r PaymentRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("PUT "+r.baseEndPoint+"/{id}", r.controller.Update)
}
