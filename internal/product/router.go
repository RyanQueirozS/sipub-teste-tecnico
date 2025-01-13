package product

import (
	"net/http"
	"sipub-test/internal"
)

type ProductRouter struct {
	baseEndPoint string
	controller   internal.IController
}

// Ideally there should be a param to change the controller when needed, but
// this will make the code bigger and I will not change the controller nor the
// repository
func NewProductRouter() ProductRouter {
	router := ProductRouter{
		controller: &ProductController{
			repository: &InMemoryProductRepository{
				products: make(map[string]ProductModel),
			},
		},
	}
	return router
}

func (r ProductRouter) Init(mux *http.ServeMux) {
	r.baseEndPoint = "/products"

	r.create(mux)
	r.getAll(mux)
	r.getOne(mux)
	r.deleteAll(mux)
	r.deleteOne(mux)
	r.update(mux)
}

func (r ProductRouter) create(mux *http.ServeMux) {
	// TODO
	mux.HandleFunc("POST "+r.baseEndPoint, r.controller.Create)
}

func (r ProductRouter) getAll(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint, r.controller.GetAll)
}

func (r ProductRouter) getOne(mux *http.ServeMux) {
	mux.HandleFunc("GET "+r.baseEndPoint+"/{id}", r.controller.GetOne)
}

func (r ProductRouter) deleteAll(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint, r.controller.DeleteAll)
}

func (r ProductRouter) deleteOne(mux *http.ServeMux) {
	mux.HandleFunc("DELETE "+r.baseEndPoint+"/{id}", r.controller.DeleteOne)
}

func (r ProductRouter) update(mux *http.ServeMux) {
	mux.HandleFunc("PUT "+r.baseEndPoint+"/{id}", r.controller.Update)
}
