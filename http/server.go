package http

import (
	"fmt"
	nethttp "net/http"
	"os"
)

func RunServer() error {
	port := ":8081"

	if os.Getenv("HTTP_PORT") != "" {
		port = os.Getenv("HTTP_PORT")
	}

	fmt.Println(port)
	router := CreateRouter()
	nethttp.Handle(string('/'), router)
	fmt.Println("Server is listening...")
	if err := nethttp.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
