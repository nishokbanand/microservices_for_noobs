package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (prod *Product) GetRequest(rw http.ResponseWriter, r *http.Request) {
	curr := r.URL.Query().Get("currency")
	lp, err := prod.p.GetProducts(curr)
	if err != nil {
		http.Error(rw, "cannot get prod", http.StatusInternalServerError)
		return
	}
	//we use NewEncoder instead of marshal to avoid having to buffer the output to an in memory slice of bytes
	// d, err :=json.Marshal(lp)
	err = lp.ToJSON(rw)
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
	curr := r.URL.Query().Get("currency")
	prod, err := p.p.GetProductByID(id, curr)
	if err != nil {
		http.Error(rw, "cannot get product", http.StatusBadRequest)
		return
	}
	//we use NewEncoder instead of marshal to avoid having to buffer the output to an in memory slice of bytes
	// d, err :=json.Marshal(lp)
	err = prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal", http.StatusInternalServerError)
		return
	}
}
