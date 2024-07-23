package handler

import (
	"log"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

type KeyProduct struct{}
