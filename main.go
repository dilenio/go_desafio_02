package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Address struct {
	Cep string `json:"cep"`
	State string `json:"state,omitempty"`
	Estado string `json:"uf,omitempty"`
	City string `json:"city,omitempty"`
	Cidade string `json:"localidade,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	Bairro string `json:"bairro,omitempty"`
	Street string `json:"street,omitempty"`
	Rua string `json:"logradouro,omitempty"`
}

func FetchAddress(apiURL string, cep string, ch chan<- *Address) {
	url := apiURL + cep
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro na requisição para %s: %s\n", apiURL, err)
		ch <- nil
		return
	}
	defer resp.Body.Close()
	var address Address
	err = json.NewDecoder(resp.Body).Decode(&address)
	if err != nil {
		fmt.Printf("Erro ao decodificar resposta da %s: %s\n", apiURL, err)
		ch <- nil
		return
	}
	ch <- &address
}

func main() {
	cep := "45208643"

	ch1 := make(chan *Address)
	ch2 := make(chan *Address)

	go FetchAddress("https://brasilapi.com.br/api/cep/v1/", cep, ch1)
	go FetchAddress("http://viacep.com.br/ws/", cep+"/json/", ch2)

	select {
	case address1 := <-ch1:
		printResult("https://brasilapi.com.br", address1)
	case address2 := <-ch2:
		printResult("http://viacep.com.br", address2)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout excedido. Nenhuma resposta recebida dentro do tempo limite.")
	}
}

func printResult(apiURL string, address *Address) {
	if apiURL == "https://brasilapi.com.br" {
		fmt.Printf("Resultado da API %s:\n", apiURL)
		fmt.Printf("CEP: %s\n", address.Cep)
		fmt.Printf("Rua: %s\n", address.Street)
		fmt.Printf("Bairro: %s\n", address.Neighborhood)
		fmt.Printf("Cidade: %s\n", address.City)
		fmt.Printf("UF: %s\n", address.State)
	} else {
		fmt.Printf("Resultado da API %s:\n", apiURL)
		fmt.Printf("CEP: %s\n", address.Cep)
		fmt.Printf("Rua: %s\n", address.Rua)
		fmt.Printf("Bairro: %s\n", address.Bairro)
		fmt.Printf("Cidade: %s\n", address.Cidade)
		fmt.Printf("UF: %s\n", address.Estado)
	}
}