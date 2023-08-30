package http

import (
	"github.com/gorilla/mux"
	nethttp "net/http"
)

type Route struct {
	Handler func(nethttp.ResponseWriter, *nethttp.Request)
	Method  string
}

func NewCourierRoute(routes map[string]Route, router *mux.Router) *mux.Router {
	for url, route := range routes {
		router.HandleFunc(url, route.Handler).Methods(route.Method)
	}

	return router
}
