package courierService

import (
	"fmt"
	"net/http"
	"os"
)

func RunServer() {
	port := ":8081"

	if os.Getenv("HTTP_PORT") != "" {
		port = os.Getenv("HTTP_PORT")
	}

	fmt.Println(port)
	router := CreateRouters()
	http.Handle(string('/'), router)
	fmt.Println("Server is listening...")
	if err := http.ListenAndServe(port, nil); err != nil {
	    return err
	}
}
