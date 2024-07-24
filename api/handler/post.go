package handler

import (
	"net/http"

	"github.com/nishokbanand/microservices/data"
)

func (p *Product) PostRequest(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
	p.l.Printf("Added: %#v", prod)
}
