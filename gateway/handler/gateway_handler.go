package handler

import (
	"fmt"
	nethttp "net/http"
	"net/http/httputil"
	neturl "net/url"
)

// NewGateWayProxyHandler create proxy for redirect on services use host with port redirect from config routes
func NewGateWayProxyHandler(hostRedirect string) func(w nethttp.ResponseWriter, r *nethttp.Request) {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		urlService, err := neturl.Parse(hostRedirect)
		fmt.Println(urlService)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)

			return
		}
		proxy := httputil.NewSingleHostReverseProxy(urlService)
		proxy.ServeHTTP(w, r)
		fmt.Print(r.RequestURI)
	}
}
