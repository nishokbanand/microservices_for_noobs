package data

import (
	"encoding/json"
	"io"
	"time"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Sku         string  `json:"sku"`
	Created_on  string  `json:"-"`
	Updated_on  string  `json:"-"`
	Deleted_on  string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(wr io.Writer) error {
	encoder := json.NewEncoder(wr)
	err := encoder.Encode(p)
	return err
}

func GetProducts() Products {
	return products
}

var products = []*Product{
	{
		ID:          "1",
		Name:        "Latte",
		Description: "sweet",
		Price:       4.50,
		Sku:         "abc123",
		Created_on:  time.Now().UTC().String(),
		Updated_on:  time.Now().UTC().String(),
	},
	{
		ID:          "2",
		Name:        "Espresso",
		Description: "Dont drink",
		Price:       2.50,
		Sku:         "123abc",
		Created_on:  time.Now().UTC().String(),
		Updated_on:  time.Now().UTC().String(),
	},
}
