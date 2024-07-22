package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

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
	} else if r.Method == http.MethodPost {
		p.postRequest(rw, r)
		return
	} else if r.Method == http.MethodPut {
		p.putRequest(rw, r)
		return
	}
	//catch all other methods
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getRequest(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	//we use NewEncoder instead of marshal to avoid having to buffer the output to an in memory slice of bytes
	// d, err :=json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal", http.StatusInternalServerError)
	}
}

func (p *Product) postRequest(rw http.ResponseWriter, r *http.Request) {
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

func (p *Product) putRequest(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	regex := regexp.MustCompile(`/([0-9]+)`)
	group := regex.FindAllStringSubmatch(path, -1)
	//FindAllStringSubmatch returns a slice of slices where the outerslice represents the matches by order
	//and the inner slice represents the exact match followed by capture groups
	//for eg : here the regex matches /[0-9]+ , when the path is usr/product/1234/details/456
	//the group becomes [[/1234,1234],[/456,456]]

	if group == nil {
		http.Error(rw, "Error in URI", http.StatusBadRequest)
	}
	fmt.Println(group[0])
	fmt.Println(group[0][1])

	//bounce off if the uri has more than one id
	if len(group[0]) != 2 {
		http.Error(rw, "Error in URI", http.StatusBadRequest)
	}
	//get the exact match
	idString := group[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "Error in URI", http.StatusBadRequest)
	}
	prod := &data.Product{}
	err = prod.FromJSON(r.Body)
	prod.ID = id
	if err != nil {
		p.l.Fatal(err)
		http.Error(rw, "Unable to Unmarshall the request", http.StatusBadRequest)
	}
	data.PutProduct(prod)
}
