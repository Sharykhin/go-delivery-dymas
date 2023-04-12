package main

import (
	"github.com/Sharykhin/go-delivery-dymas/location/http"
	"log"
)

func main() {
	if err := http.RunServer(); err != nil {
		log.Printf("failed to run http server: %v", err)
	}
}
