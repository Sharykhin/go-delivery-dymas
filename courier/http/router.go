package http

import (
	http "github.com/Sharykhin/go-delivery-dymas/courier/http/handler"
	"github.com/gorilla/mux"
)

type RouteCourier struct {
	url string
}

func NewCourierCreateRoute() *RouteCourier {
	return &RouteCourier{
		url: "/couriers",
	}
}

func (r *RouteCourier) NewCourierCreateRoute(courierHandler *http.CourierHandler, router *mux.Router) *mux.Router {
	router.HandleFunc(r.url, courierHandler.HandlerCourierCreate).Methods("POST")

	return router
}
