package main

import (
	"log"
	"net"
	"os"

	protos "github.com/nishokbanand/learngrpc/protos/currency"
	"github.com/nishokbanand/learngrpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := log.New(os.Stdout, "grpc logger>>", log.Default().Flags())
	gs := grpc.NewServer()
	reflection.Register(gs)
	cs := server.NewCurrencyService(log, protos.UnimplementedCurrencyServer{})
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
