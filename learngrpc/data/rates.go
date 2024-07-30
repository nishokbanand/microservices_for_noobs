package data

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

type ExchangeRates struct {
	l     *log.Logger
	rates map[string]float64
}

func NewRates(l *log.Logger) (*ExchangeRates, error) {
	ex := &ExchangeRates{l, map[string]float64{}}
	ex.getRates()
	return ex, nil
}

func (ex *ExchangeRates) GetRate(base, destination string) (float64, error) {
	bv, ok := ex.rates[base]
	if !ok {
		return 0, fmt.Errorf("could not find the currency %s", base)
	}
	dv, ok := ex.rates[destination]
	if !ok {
		return 0, fmt.Errorf("could not find the currency %s", destination)
	}
	return dv / bv, nil
}

var URI = "" //supposed to be european exchange xml api but no longer available

func getDataFromFile() ([]byte, error) {
	filePath := getPathOfFile()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getPathOfFile() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "eurofxref-daily.xml")
	return d
}

func (ex *ExchangeRates) getRates() error {
	// resp, err := http.Get(URI)
	resp, err := getDataFromFile()
	if err != nil {
		return err
	}
	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("error in getting the api %d", resp.StatusCode)
	// }
	// defer resp.Body.Close()
	// decoder := xml.NewDecoder(resp.Body)
	reader := bytes.NewReader(resp)
	decoder := xml.NewDecoder(reader)

	Cubes := &Cubes{}
	decoder.Decode(&Cubes)
	for _, rate := range Cubes.CubeData {
		r, err := strconv.ParseFloat(rate.Rate, 64)
		if err != nil {
			return err
		}
		ex.rates[rate.Currency] = r
	}
	ex.rates["EUR"] = 1
	return nil
}

func (ex *ExchangeRates) MontiorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				for k, v := range ex.rates {
					change := rand.Float64() / 10
					dir := rand.Intn(1)
					if dir == 0 {
						change = 1 - change
					} else {
						change = 1 + change
					}
					ex.rates[k] = v * change
				}
				ret <- struct{}{}
			}
		}
	}()
	return ret
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
