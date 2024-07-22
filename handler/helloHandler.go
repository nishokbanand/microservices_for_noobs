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
	//r.Body implements io.Reader so io.ReadAll can read it
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error in HelloHandler", http.StatusBadRequest)
	}
	// you get back a byte array, so use a formatter to give it string
	fmt.Printf("Hello : %s", d)
}
