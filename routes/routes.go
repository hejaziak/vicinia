package routes

import (
	"net/http"
	handlers "vicinia/handlers"
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
		handlers.IndexHandler,
	},
	Route{
		"Welcome",
		"GET",
		"/welcome",
		handlers.WelcomeHandler,
	},
	Route{
		"chat",
		"POST",
		"/chat",
		handlers.ChatHandler,
	},
}
