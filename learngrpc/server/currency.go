package server

import (
	"context"
	"log"

	protos "github.com/nishokbanand/learngrpc/protos/currency"
)

type Currency struct {
	l *log.Logger
	protos.UnimplementedCurrencyServer
}

func NewCurrencyService(l *log.Logger, uc protos.UnimplementedCurrencyServer) *Currency {
	return &Currency{l, uc}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.l.Println("base", rr.GetBase(), "destination", rr.GetDestination())
	return &protos.RateResponse{
		Rate: 0.5,
	}, nil
}
