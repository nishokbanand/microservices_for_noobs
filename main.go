package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/nishokbanand/microservices/handler"
)

func main() {
	logger := log.New(os.Stdout, "logger >>", log.Default().Flags())
	sm := mux.NewRouter()

	//Create a handler with method ServeHTTP
	ph := handler.NewProduct(logger)
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetRequest)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.PostRequest)
	postRouter.Use(ph.MiddleWareValidateProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.PutRequest)
	putRouter.Use(ph.MiddleWareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteRequest)

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
