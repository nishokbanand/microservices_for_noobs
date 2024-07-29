package main

import (
	"log"
	"net"
	"os"

	"github.com/nishokbanand/learngrpc/data"
	protos "github.com/nishokbanand/learngrpc/protos/currency"
	"github.com/nishokbanand/learngrpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := log.New(os.Stdout, "grpc logger>>", log.Default().Flags())
	ex, err := data.NewRates(log)
	if err != nil {
		log.Println("error in getting rates", err)
	}
	gs := grpc.NewServer()
	reflection.Register(gs)
	cs := server.NewCurrencyService(log, ex, protos.UnimplementedCurrencyServer{})
	protos.RegisterCurrencyServer(gs, cs)
	netListener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Println("error in creating net listener", err)
	}
	err = gs.Serve(netListener)
	if err != nil {
		log.Println("error in grpc server", err)
	}
}
