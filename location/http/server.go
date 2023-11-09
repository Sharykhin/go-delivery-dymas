package http

import (
	"context"
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// RunServer  runs courier http server.
func RunServer(ctx context.Context, router *mux.Router, port string) {
	fmt.Println(port)
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
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	fmt.Println("Stop Server signal")
}
