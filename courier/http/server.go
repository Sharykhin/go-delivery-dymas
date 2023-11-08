package http

import (
	"fmt"
	nethttp "net/http"

	"github.com/gorilla/mux"
)

// RunServer runs http server
func RunServer(router *mux.Router, port string) error {
	fmt.Println(port)
	nethttp.Handle(string('/'), router)
	fmt.Println("Server is listening...")
	if err := nethttp.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
