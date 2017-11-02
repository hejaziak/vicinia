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
