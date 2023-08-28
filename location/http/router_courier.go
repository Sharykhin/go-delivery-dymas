package http

import (
	"fmt"
	http "github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/gorilla/mux"
)

type Router struct {
	url        string
	uuidRegexp string
}

func NewRouter() *Router {
	return &Router{
		url:        "/courier/{courier_id:%s}/location",
		uuidRegexp: "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
	}
}

func (r *Router) NewRouter(locationHandler *http.LocationHandler, router *mux.Router) *mux.Router {
	r.url = fmt.Sprintf(r.url, r.uuidRegexp)
	router.HandleFunc(r.url, locationHandler.HandlerCouriersLocation).Methods("POST")

	return router
}
