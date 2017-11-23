package routes

import (
	"net/http"
	handlers "vicinia/handlers"
	middleware "vicinia/middleware"

	"github.com/gorilla/mux"
)

//NewRouter : creates a new router and returns it
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.NotFoundHandler = http.HandlerFunc(handlers.MiscHandler)
	return router
}
