package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/nishokbanand/learngrpc/data"
	protos "github.com/nishokbanand/learngrpc/protos/currency"
)

type Currency struct {
	l  *log.Logger
	ex *data.ExchangeRates
	protos.UnimplementedCurrencyServer
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrencyService(l *log.Logger, ex *data.ExchangeRates, uc protos.UnimplementedCurrencyServer) *Currency {
	c := &Currency{l, ex, uc, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()
	return c
}

func (c *Currency) handleUpdates() {
	ru := c.ex.MontiorRates(time.Second * 10)
	for range ru {
		c.l.Println("GOT updates")
		//go through each subs
		for k, v := range c.subscriptions {
			//go through each RateRequest
			for _, rr := range v {
				r, err := c.ex.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.l.Println("error in getting the rates", err)
				}
				err = k.Send(&protos.RateResponse{
					Base:        rr.Base,
					Destination: rr.Destination,
					Rate:        r,
				})
				if err != nil {
					c.l.Println("error in sending the response", err)
				}
			}
		}
	}
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

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	for {
		rr, err := src.Recv()
		c.l.Println("in")
		if err == io.EOF {
			c.l.Println("client closed the connection", err)
			break
		}
		if err != nil {
			c.l.Println("error in receiving msg", err)
			break
		}
		c.l.Println("request base", rr.GetBase().String(), "destination", rr.GetDestination().String())
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}
	return nil
}
