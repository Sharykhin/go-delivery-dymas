package handler

import (
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	"net/http/httputil"
	neturl "net/url"

	pkghttp "github.com/Sharykhin/go-delivery-dymas/pkg/http"
)

// NewGateWayProxyHandler create proxy for redirect on services use host with port redirect from config routes
func NewGateWayProxyHandler(hostRedirect string) func(w nethttp.ResponseWriter, r *nethttp.Request) {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		urlService, err := neturl.Parse(hostRedirect)
		fmt.Println(urlService)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			err := json.NewEncoder(w).Encode(&pkghttp.ResponseMessage{
				Status:  "Error",
				Message: "Incorrect config host",
			})

			if err != nil {
				log.Printf("failed to encode json response: %v\n", err)
			}

			return
		}
		proxy := httputil.NewSingleHostReverseProxy(urlService)
		proxy.ServeHTTP(w, r)
		fmt.Print(r.RequestURI)
	}
}
