package http

import (
	"fmt"
	http "github.com/Sharykhin/go-delivery-dymas/location/http/handler"
	"github.com/gorilla/mux"
)

type RouteCourierCreate struct {
	url string
}

func NewCourierCreateRoute() *RouteCourierCreate {
	return &RouteCourierCreate{
		url: "/couriers",
	}
}

func (r *RouteCourierCreate) NewCourierCreateRoute(locationHandler *http.CourierCreateHandler, router *mux.Router) *mux.Router {
	router.HandleFunc(r.url, locationHandler.HandlerCourierCreate).Methods("POST")

	return router
}
