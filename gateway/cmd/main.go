package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/Sharykhin/go-delivery-dymas/gateway/env"
	"github.com/Sharykhin/go-delivery-dymas/gateway/middleware"
	"github.com/Sharykhin/go-delivery-dymas/gateway/route"
	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

func main() {
	config, err := env.GetConfig()
	if err != nil {
		log.Printf("failed to parse variable env: %v\n", err)
		return
	}
	middlewares := map[string]func(paramName string) func(next http.Handler) http.Handler{
		"uuid": middleware.UuidMiddleware,
	}
	routes, err := route.CreateServicesRoutesFromConfig(middlewares)
	if err != nil {
		log.Printf("failed to parse routes config: %v\n", err)
		return
	}

	requiresMiddlewares := []func(next http.Handler) http.Handler{pkghttp.GetRequestID, pkghttp.CreateReqIDMiddleware}
	router := pkghttp.NewRoute(routes, mux.NewRouter(), requiresMiddlewares)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pkghttp.RunServer(ctx, router, ":"+config.PortServerGateway)
}
