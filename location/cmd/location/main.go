package main

import (
	"fmt"
	"github.com/Sharykhin/go-delivery-dymas/location/http"
)

func main() {
	err := http.RunServer()

	if err != nil {
		fmt.Println(err)
	}
}
