package internal

import "net/http"

func RouterInitializeAll(mux *http.ServeMux, routers []IRouter) {
	for i := 0; i < len(routers); i++ {
		routers[i].Init(mux)
	}
}
