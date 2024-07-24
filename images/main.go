package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/nishokbanand/imageHandler/handler"
)

func main() {
	sm := mux.NewRouter()

	l := log.New(os.Stdout, "image-logger >>", log.Default().Flags())

	fh := handler.NewFiles(l)
	getImgReq := sm.Methods(http.MethodGet).Subrouter()
	getImgReq.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}", fh.GetRequest())

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
		IdleTimeout:  time.Second * 120,
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	fmt.Println("Shutting down", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)
}
