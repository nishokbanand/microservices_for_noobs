package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nishokbanand/microservices/handler"
)

func main() {
	logger := log.New(os.Stdout, "logger >>", log.Default().Flags())
	hh := handler.NewHello(logger)
	gh := handler.NewGoodBye(logger)
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
		IdleTimeout:  time.Second * 120,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)
	sig := <-sigchan
	logger.Println("recieve terminate", sig)
	//gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)
}
