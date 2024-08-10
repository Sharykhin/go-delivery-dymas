package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"gateway/env"
	"gateway/middleware"
	"gateway/route"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

func main() {
	middlewares := map[string]func(paramName string) func(next http.Handler) http.Handler{
		"uuid": middleware.UuidMiddleware,
	}
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	routes, err := route.CreateServicesRoutesFromConfig(middlewares)
	if err != nil {
		log.Printf("failed to parse routes config: %v\n", err)
		return
	}
	router := pkghttp.NewRoute(routes, mux.NewRouter())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pkghttp.RunServer(ctx, router, ":"+config.PortServerCourier)
}
