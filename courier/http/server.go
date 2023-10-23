package http

import (
	"fmt"
	"github.com/gorilla/mux"
	nethttp "net/http"
)

func RunServer(router *mux.Router, port string) error {
	fmt.Println(port)
	nethttp.Handle(string('/'), router)
	fmt.Println("Server is listening...")
	if err := nethttp.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
