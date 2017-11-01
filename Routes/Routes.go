package routes

import (
	h "Vicinia/Handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		h.Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		h.TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		h.TodoShow,
	},
	Route{
		"Welcome",
		"GET",
		"/welcome",
		h.WelcomeHandler,
	},
	Route{
		"chat",
		"GET",
		"/chat",
		h.ListHandler,
	},
}
