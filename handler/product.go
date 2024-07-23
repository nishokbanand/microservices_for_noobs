package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nishokbanand/microservices/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}
func (p *Product) GetRequest(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	//we use NewEncoder instead of marshal to avoid having to buffer the output to an in memory slice of bytes
	// d, err :=json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal", http.StatusInternalServerError)
	}
}

func (p *Product) PostRequest(rw http.ResponseWriter, r *http.Request) {
	//We use NewDecoder instead of unmarshal
	// d, _ := io.ReadAll(r.Body)
	// json.Unmarshal(d, &data.Product{})
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to Unmarshall the request", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Added: %#v", prod)
}

type KeyProduct struct{}

func (p *Product) PutRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Fatal(err)
		http.Error(rw, "Unable to get ID", http.StatusBadRequest)
	}
	println("here")
	fmt.Println(r.Context().Value(KeyProduct{}))
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	prod.ID = id
	err = data.PutProduct(prod)
	if err != nil {
		p.l.Fatal(err)
		http.Error(rw, "Unable to get the product", http.StatusBadRequest)
	}
}

func (p *Product) MiddleWareFromJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Fatal(err)
			http.Error(rw, "Unable to Unmarshall the request", http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
