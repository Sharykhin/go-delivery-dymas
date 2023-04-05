package main

import (
	"fmt"
	"http/http"
)

func main() {
	err := http.RunServer()

	if err != nil {
		fmt.Println(err)
	}
}
