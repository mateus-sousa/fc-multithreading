package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	postalCode := os.Args[1]
	c1 := make(chan []byte)
	c2 := make(chan []byte)
	go getAddressByViaCep(postalCode, c1)
	go getAddressByBrasilAPI(postalCode, c2)

	select {
	case address := <-c1:
		fmt.Println(string(address))
	case address := <-c2:
		fmt.Println(string(address))
	case <-time.After(time.Second * 1):
		log.Fatal("request to get address by postalcode timeout")
	}
}

func getAddressByViaCep(postalCode string, c1 chan []byte) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", postalCode), nil)
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
