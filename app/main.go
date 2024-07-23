package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	mux.GET("/youtube/channel/stats", getChannelStats())

	return mux
}

func getChannelStats() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello, World!"))
	}
}

func main() {
	println("Starting the server...")
	println("Listening on port 1010")

	srv := &http.Server{
		Addr:    ":1010",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		println("Shutting down the server...")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("fatal http server failed to shutdown: %v", err)
		}

		log.Println("Server shutdown successfully")

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}
	<-idleConnsClosed
	log.Println("Server shutdown completed")
}
