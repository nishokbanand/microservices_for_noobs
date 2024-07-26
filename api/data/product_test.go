package data

import (
	"testing"
)

func TestValidation(t *testing.T) {
	prod := &Product{
		Name:  "tea",
		Price: 1,
		Sku:   "abc-def-ghi",
	}
	err := prod.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
