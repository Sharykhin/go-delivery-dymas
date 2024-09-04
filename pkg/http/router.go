package http

import (
	"context"
	"fmt"
	"net/http"
	nethttp "net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const RequestIDKeyHeader = "X-Request-ID"
const RequestIDKeyContextValue = "requestID"

// Route handles different path routes
type Route struct {
	Handler     func(nethttp.ResponseWriter, *nethttp.Request)
	Methods     []string
	Middlewares []func(next http.Handler) http.Handler
}

// NewRoute creates for handling different path routes
func NewRoute(routes map[string]Route, router *mux.Router, additionalMiddlewares []func(next http.Handler) http.Handler) *mux.Router {
	for url, route := range routes {
		if route.Middlewares != nil && additionalMiddlewares != nil {
			route.Middlewares = append(route.Middlewares, additionalMiddlewares...)
		} else if route.Middlewares == nil && additionalMiddlewares != nil {
			route.Middlewares = additionalMiddlewares
		}

		if route.Middlewares != nil {
			handle := prepareMiddleware(http.HandlerFunc(route.Handler), route.Middlewares...)
			router.Handle(url, handle).Methods(route.Methods...)
		} else {
			router.HandleFunc(url, route.Handler).Methods(route.Methods...)
		}
	}

	return router
}

func prepareMiddleware(handler http.Handler, Middleware ...func(next http.Handler) http.Handler) http.Handler {
	for i := len(Middleware); i > 0; i-- {
		handler = Middleware[i-1](handler)
	}
	return handler
}

// CreateReqIDMiddleware middleware set requestID that has uuid format requestId needs for logging requests
func CreateReqIDMiddleware(next http.Handler) http.Handler {
	fmt.Println("CreateReqIDMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		requestID := id.String()
		fmt.Println(requestID)
		r.Header.Set(RequestIDKeyHeader, requestID)
		next.ServeHTTP(w, r)
	})
}

// GetRequestID get requestID that has uuid format requestId needs for logging requests
func GetRequestID(next http.Handler) http.Handler {
	fmt.Println("ContextReqIDMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ContextReqIDMiddleware")
		reqID := r.Header.Get(RequestIDKeyHeader)
		ctx := r.Context()
		fmt.Println(reqID)
		ctx = context.WithValue(ctx, RequestIDKeyContextValue, reqID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
