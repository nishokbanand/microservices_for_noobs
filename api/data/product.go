package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	protos "github.com/nishokbanand/learngrpc/protos/currency"
)

type Product struct {
	ID          int     `json:"id" `
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	Sku         string  `json:"sku" validate:"required,sku"`
	Created_on  string  `json:"-"`
	Updated_on  string  `json:"-"`
	Deleted_on  string  `json:"-"`
}

func (p *Product) Validate() error {
	validator := validator.New()
	validator.RegisterValidation("sku", skuValidation)
	return validator.Struct(p)
}

func skuValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]`)
	matches := regex.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}
	return true
}

type Products []*Product

func (p *Product) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(p)
	return err
}

func (p *Products) ToJSON(wr io.Writer) error {
	encoder := json.NewEncoder(wr)
	err := encoder.Encode(p)
	return err
}

type ProductsDB struct {
	l *log.Logger
	c protos.CurrencyClient
}

func NewProductsDB(l *log.Logger, c protos.CurrencyClient) *ProductsDB {
	return &ProductsDB{l, c}
}

func (p *ProductsDB) GetProducts(dest_curr string) (Products, error) {
	if dest_curr == "" {
		return products, nil
	}
	rate, err := p.getRate(dest_curr)
	if err != nil {
		p.l.Println(err)
		return nil, err
	}
	prod := Products{}
	for _, p := range products {
		np := *p
		np.Price = np.Price * rate
		prod = append(prod, &np)
	}
	return prod, nil
}

func (p *ProductsDB) GetProductByID(id int, dest_curr string) (*Product, error) {
	if dest_curr == "" {
		return products[id], nil
	}
	rate, err := p.getRate(dest_curr)
	if err != nil {
		p.l.Println(err)
		return nil, err
	}
	np := *products[id]
	np.Price = np.Price * rate
	return &np, nil
}

func AddProduct(p *Product) {
	id := getNextId()
	p.ID = id
	products = append(products, p)
}

func PutProduct(p *Product) error {
	idx, err := GetProduct(p.ID)
	if err != nil {
		return err
	}
	products[idx] = p
	return nil
}
func (p *Product) ToJSON(wr io.Writer) error {
	encoder := json.NewEncoder(wr)
	err := encoder.Encode(p)
	return err
}

func DeleteProduct(id int) error {
	idx, err := GetProduct(id)
	if err != nil {
		return err
	}
	products = append(products[:idx], products[idx+1:]...)
	return nil
}

func FindProduct(id int) (*Product, error) {
	idx, err := GetProduct(id)
	if err != nil {
		return nil, err
	}
	return products[idx], nil
}

func GetProduct(id int) (int, error) {
	for idx, value := range products {
		if value.ID == id {
			return idx, nil
		}
	}
	return -1, ProductNotFound
}

var ProductNotFound = fmt.Errorf("Product not found")

func getNextId() int {
	id := products[len(products)-1].ID
	return id + 1
}
func (p *ProductsDB) getRate(dest_curr string) (float64, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[dest_curr]),
	}
	fmt.Println(rr.Base.String())
	fmt.Println(rr.Destination.String())
	resp, err := p.c.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println(err)
		return 0, err
	}
	p.l.Println("Rate is", resp.Rate)
	return resp.Rate, nil
}

var products = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "sweet",
		Price:       4.50,
		Sku:         "abc-def-ghi",
		Created_on:  time.Now().UTC().String(),
		Updated_on:  time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Dont drink",
		Price:       2.50,
		Sku:         "ghi-def-abc",
		Created_on:  time.Now().UTC().String(),
		Updated_on:  time.Now().UTC().String(),
	},
}
