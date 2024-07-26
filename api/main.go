package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	protos "github.com/nishokbanand/learngrpc/protos/currency"
	"github.com/nishokbanand/microservices/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := log.New(os.Stdout, "logger >>", log.Default().Flags())
	sm := mux.NewRouter()
	conn, err := grpc.NewClient("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		logger.Fatal("cannot make a connection")
	}
	cc := protos.NewCurrencyClient(conn)

	//Create a handler with method ServeHTTP
	ph := handler.NewProduct(logger, cc)
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetRequest)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.ListOneProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.PostRequest)
	postRouter.Use(ph.MiddleWareValidateProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.PutRequest)
	putRouter.Use(ph.MiddleWareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteRequest)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"locahost:9090"}))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
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
