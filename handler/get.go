package handler

import (
	"net/http"

	"github.com/nishokbanand/microservices/data"
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
