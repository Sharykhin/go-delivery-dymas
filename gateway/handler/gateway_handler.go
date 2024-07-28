package handler

import (
	"fmt"
	nethttp "net/http"
	"net/http/httputil"
	neturl "net/url"
	"regexp"
)

type GateWayProxyHandler struct {
	routes map[string]string
}

func NewGateWayProxyHandler(routes map[string]string) *GateWayProxyHandler {
	return &GateWayProxyHandler{
		routes: routes,
	}
}

func (gph *GateWayProxyHandler) RequestServiceHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
	for url, host := range gph.routes {
		isMatch, _ := regexp.MatchString("^"+url+"$", r.RequestURI)
		if isMatch {
			urlService, err := neturl.Parse(host)
			fmt.Println(urlService)
			if err != nil {
				w.WriteHeader(nethttp.StatusBadRequest)

				return
			}
			proxy := httputil.NewSingleHostReverseProxy(urlService)
			proxy.ServeHTTP(w, r)
		}
	}
	fmt.Print(r.RequestURI)
}
