package data

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	// "net/http"
	"os"
	"strconv"
)

type ExchangeRates struct {
	l     *log.Logger
	rates map[string]float64
}

func NewRates(l *log.Logger) (*ExchangeRates, error) {
	getDataFromFile()
	ex := &ExchangeRates{l, map[string]float64{}}
	ex.getRates()
	return ex, nil
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
	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
