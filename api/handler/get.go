package handler

import (
	"context"
	"github.com/gorilla/mux"
	protos "github.com/nishokbanand/learngrpc/protos/currency"
	"github.com/nishokbanand/microservices/data"
	"net/http"
	"strconv"
)

func (p *Product) GetRequest(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	//we use NewEncoder instead of marshal to avoid having to buffer the output to an in memory slice of bytes
	// d, err :=json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal", http.StatusInternalServerError)
		return
	}
}

func (p *Product) ListOneProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "cannot get id", http.StatusInternalServerError)
		return
	}
	prod, err := data.FindProduct(id)
	if err != nil {
		http.Error(rw, "cannot get product", http.StatusBadRequest)
		return
	}
	//we use NewEncoder instead of marshal to avoid having to buffer the output to an in memory slice of bytes
	// d, err :=json.Marshal(lp)
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value["GBP"]),
	}
	resp, err := p.c.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println(err)
		http.Error(rw, "cannot get conversion rate", http.StatusInternalServerError)
		return
	}
	p.l.Println("Rate is", resp.Rate)
	prod.Price = resp.Rate * prod.Price
	err = prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal", http.StatusInternalServerError)
		return
	}
}
