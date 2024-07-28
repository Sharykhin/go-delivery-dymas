package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gateway/handler"
	"gateway/route"
	"github.com/gorilla/mux"

	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

func main() {
	sevicesroutes := route.CreateServicesRoutes()
	gateWayProxyHandler := handler.NewGateWayProxyHandler(sevicesroutes)
	routes := map[string]pkghttp.Route{"/": {
		Handler:    gateWayProxyHandler.RequestServiceHandler,
		Method:     []string{"GET", "POST", "PUT", "PATCH"},
		PathPrefix: "/",
	},
	}
	router := pkghttp.NewRoute(routes, mux.NewRouter())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pkghttp.RunServer(ctx, router, ":"+"8080")
}
