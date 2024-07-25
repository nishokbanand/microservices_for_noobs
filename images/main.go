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
	"github.com/nishokbanand/imageHandler/files"
	"github.com/nishokbanand/imageHandler/handler"
)

func main() {
	sm := mux.NewRouter()

	l := log.New(os.Stdout, "image-logger >>", log.Default().Flags())

	basePath := "./imagestore"
	store, err := files.NewLocal(basePath, 1024*1000*5)
	if err != nil {
		l.Fatalf("cannot create a local storage: %v", err)
	}
	fh := handler.NewFiles(l, store)
	postImgReq := sm.Methods(http.MethodPost).Subrouter()
	postImgReq.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)

	getImgReq := sm.Methods(http.MethodGet).Subrouter()
	getImgReq.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/images/", http.FileServer(http.Dir(basePath))))
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
			l.Println("Error in running server", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	fmt.Println("Shutting down", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)
}
