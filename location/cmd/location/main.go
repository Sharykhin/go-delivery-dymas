package main

import (
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
)

func main() {
	if err := http.RunServer(); err != nil {
	    log.Printf("failed to run http server: %v", err)
	}

	if err != nil {
		fmt.Println(err)
	}
}
