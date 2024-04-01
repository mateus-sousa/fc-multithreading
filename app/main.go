package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	c1 := make(chan []byte)
	c2 := make(chan []byte)
	go getAddressByViaCep("73402574", c1)
	go getAddressByBrasilAPI("73402574", c2)
}

func getAddressByViaCep(cepNumber string, c1 chan []byte) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cepNumber), nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(errors.New("error to request cep"))
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	c1 <- body
}

func getAddressByBrasilAPI(cepNumber string, c2 chan []byte) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v2/%s", cepNumber), nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(errors.New("error to request cep"))
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	c2 <- body
}
