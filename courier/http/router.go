package http

import (
	http "github.com/Sharykhin/go-delivery-dymas/courier/http/handler"
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

func (r *RouteCourierCreate) NewCourierCreateRoute(courierHandler *http.CourierHandler, router *mux.Router) *mux.Router {
	router.HandleFunc(r.url, courierHandler.HandlerCourierCreate).Methods("POST")

	return router
}
