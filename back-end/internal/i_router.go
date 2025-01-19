package internal

import "net/http"

// This interface will be used by routers to initialize the routes. Each
// function (create for example), will abstract a method from it's respective
// controller.
// This interface will be used for a function (RouterInitializeAll) that
// recieves multiple routers and initializes them.
type IRouter interface {
	Init(mux *http.ServeMux)

	// These functions below should be contained in most structures, but golang
	// doesn't support it to be private (unexported method X).
	// Also, as the project grows, it's hard to garantee that all routers will
	// use these methods.

	// Should call the Create method for the respective controller
	// create(mux *http.ServeMux)

	// Should call the GetAll method for the respective controller
	// getAll(mux *http.ServeMux)

	// Should call the GetOne method for the respective controller
	// getOne(mux *http.ServeMux)

	// Should call the DeleteOne method for the respective controller
	// deleteOne(mux *http.ServeMux)

	// Should call the DeleteAll method for the respective controller
	// deleteAll(mux *http.ServeMux)

	// Should call the Update method for the respective controller
	// update(mux *http.ServeMux)
}
