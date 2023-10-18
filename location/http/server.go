package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunServer(router *mux.Router, port string) {
	fmt.Println(port)
	nethttp.Handle(string('/'), router)
	fmt.Println("Server is listening...")
	srv := &nethttp.Server{
		Addr:    port,
		Handler: nil,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	fmt.Println("Stop Server signal")
}
