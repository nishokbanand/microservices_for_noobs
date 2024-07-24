package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nishokbanand/microservices/data"
)

func (p *Product) MiddleWareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] in Unmarshall", err)
			http.Error(rw, "Unable to Unmarshall the request", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] in validation", err)
			http.Error(rw, fmt.Sprintf("Unable to Validate the request : %s", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
