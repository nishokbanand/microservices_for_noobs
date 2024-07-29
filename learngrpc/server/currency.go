package server

import (
	"context"
	"log"

	"github.com/nishokbanand/learngrpc/data"
	protos "github.com/nishokbanand/learngrpc/protos/currency"
)

type Currency struct {
	l  *log.Logger
	ex *data.ExchangeRates
	protos.UnimplementedCurrencyServer
}

func NewCurrencyService(l *log.Logger, ex *data.ExchangeRates, uc protos.UnimplementedCurrencyServer) *Currency {
	return &Currency{l, ex, uc}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.l.Println("base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.ex.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		c.l.Println(err)
	}
	return &protos.RateResponse{
		Rate: rate,
	}, nil
}
