package handler

import (
	protos "github.com/nishokbanand/learngrpc/protos/currency"
	"log"
)

type Product struct {
	l *log.Logger
	c protos.CurrencyClient
}

func NewProduct(l *log.Logger, cc protos.CurrencyClient) *Product {
	return &Product{l, cc}
}

type KeyProduct struct{}
