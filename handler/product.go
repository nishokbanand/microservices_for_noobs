package handler

import (
	"log"
	"net/http"

	"github.com/nishokbanand/microservices/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getRequest(rw, r)
		return
	}
	//catch all other methods
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getRequest(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	//we use NewEncoder instead of marshal to avoid storing the data in buffer
	// d, err :=json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal", http.StatusInternalServerError)
	}
}
