package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nishokbanand/microservices/data"
)

func (p *Product) DeleteRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Println("[ERROR] in Getting ID", err)
		http.Error(rw, "Unable to get ID", http.StatusBadRequest)
		return
	}
	println("here")
	err = data.DeleteProduct(id)
	if err != nil {
		//Fatal does os.Exit(1) after printing
		p.l.Println("[ERROR] in getting the product", err)
		http.Error(rw, "Unable to get the product", http.StatusBadRequest)
		return
	}
}
