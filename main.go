package main

import (
	"context"
	"github.com/nishokbanand/microservices/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "logger >>", log.Default().Flags())
	// hh := handler.NewHello(logger)
	gh := handler.NewGoodBye(logger)
	// we dont use the defaultServeMux (we are better)
	sm := http.NewServeMux()
	//Create a handler with method ServeHTTP
	ph := handler.NewProduct(logger)

	sm.Handle("/", ph)
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
	//use a bufferedChannel (to be safe) because if there is any line that comes after the sig:=<-sigchan, it wont be executed
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)
	// big brain move, channel waits for the message and it is done after the singal notfies
	sig := <-sigchan
	logger.Println("recieve terminate", sig)
	//gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	//just cancel this guy so that we get no warnings
	defer cancel()
	server.Shutdown(ctx)
}
