package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type GoodByeHandler struct {
	l *log.Logger
}

func NewGoodBye(l *log.Logger) *GoodByeHandler {
	return &GoodByeHandler{l: l}
}

func (h *GoodByeHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Good Bye")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error in GoodByeHandler", http.StatusBadRequest)
	}
	fmt.Printf("GoodBye : %s", d)
}
