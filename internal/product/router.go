package product

import (
	"net/http"
)

type ProductRouter struct {
	baseEndPoint string
}

func (r ProductRouter) Init(mux *http.ServeMux) {
	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteOne(mux)
	r.deleteAll(mux)
	r.update(mux)
}

func (r ProductRouter) create(mux *http.ServeMux) {
	// TODO
	mux.HandleFunc("Method "+r.baseEndPoint, nil)
}

func (r ProductRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("", nil)
}

func (r ProductRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("", nil)
}

func (r ProductRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("", nil)
}

func (r ProductRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("", nil)
}

func (r ProductRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("", nil)
}
