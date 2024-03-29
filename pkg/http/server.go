package http

import (
	"context"
	"fmt"
	"log"
	nethttp "net/http"
	"time"

	"github.com/gorilla/mux"
)

// RunServer  runs http server.
func RunServer(ctx context.Context, router *mux.Router, port string) {
	nethttp.Handle(string('/'), router)
	fmt.Println("Server is listening...")
	srv := &nethttp.Server{
		Addr:    port,
		Handler: nil,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	fmt.Println("Stop Server signal")
}
