package routes

import (
	handlers "Vicinia/Handlers"
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
		handlers.Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		handlers.TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		handlers.TodoShow,
	},
	Route{
		"Welcome",
		"GET",
		"/welcome",
		handlers.WelcomeHandler,
	},
	Route{
		"chat",
		"GET",
		"/chat",
		handlers.ListHandler,
	},
}
