package data

import (
	"fmt"
	"log"
	"testing"
)

func TestNewRates(t *testing.T) {
	ex, err := NewRates(log.Default())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Rates %#v", ex.rates)
}
