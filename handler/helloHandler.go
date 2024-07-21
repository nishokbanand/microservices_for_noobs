package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type HelloHandler struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *HelloHandler {
	return &HelloHandler{l: l}
}

func (h *HelloHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("hello")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error in HelloHandler", http.StatusBadRequest)
	}
	fmt.Printf("Hello : %s", d)
}
