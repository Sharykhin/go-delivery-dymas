package http

import (
	"github.com/gorilla/mux"
	nethttp "net/http"
)

const UuidRegexp = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"

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
