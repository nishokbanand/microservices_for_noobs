package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id" `
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
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

func GetProducts() Products {
	return products
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

func DeleteProduct(id int) error {
	idx, err := GetProduct(id)
	if err != nil {
		return err
	}
	products = append(products[:idx], products[idx+1:]...)
	return nil
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
