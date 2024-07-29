package handler

import (
	"log"

	"github.com/nishokbanand/microservices/data"
)

type Product struct {
	l *log.Logger
	p *data.ProductsDB
}

func NewProduct(l *log.Logger, p *data.ProductsDB) *Product {
	return &Product{l, p}
}

type KeyProduct struct{}
