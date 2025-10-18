package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const sample_cep = "01153000"

type BuscadorEndereco interface {
	BuscaByCEP(cep string) (Endereco, error)
}

type ViaCepAPI struct {
}

type BrasilAPI struct {
}

type Endereco struct {
	Estado string
	Cidade string
	Bairro string
	Rua    string
}

func main() {
	buscadorBrasilApi := BrasilAPI{}
	buscadorViaCepApi := ViaCepAPI{}

	chBrasilApi := make(chan Endereco)
	chViaCepApi := make(chan Endereco)

	var cep string
	if len(os.Args) > 1 {
		cep = os.Args[1]
	}
	if cep == "" {
		cep = sample_cep
		fmt.Printf("Executando programa com cep de exemplo: %s\n", cep)
	}

	go func() {
		endereco, err := buscadorBrasilApi.BuscaByCEP(cep)
		if err != nil {
			fmt.Printf("erro ao consultar o cep %s na api Brasil API\n", cep)
		}
		chBrasilApi <- endereco
	}()

	go func() {
		endereco, err := buscadorViaCepApi.BuscaByCEP(cep)
		if err != nil {
			fmt.Printf("erro ao consultar o cep %s no api Via CEP API\n", cep)
		}
		chViaCepApi <- endereco
	}()

	select {
	case endereco := <-chBrasilApi:
		fmt.Printf("Brasil API - CEP %s consultado com sucesso: %v\n", cep, endereco)
	case endereco := <-chViaCepApi:
		fmt.Printf("Via CEP API - CEP %s consultado com sucesso: %v\n", cep, endereco)
	case <-time.After(time.Second):
		fmt.Printf("Timeout")
	}

}

func (v *BrasilAPI) BuscaByCEP(cep string) (Endereco, error) {
	resp, err := http.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep))
	if err != nil {
		return Endereco{}, fmt.Errorf("erro ao consumir api %v", err)
	}

	type cepResponse struct {
		State        string `json:"state"`
		City         string `json:"city"`
		Neighborhood string `json:"neighborhood"`
		Street       string `json:"street"`
	}

	var response cepResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Endereco{}, fmt.Errorf("erro ao realizar o decode %v", err)
	}

	return Endereco{
		Estado: response.State,
		Cidade: response.City,
		Bairro: response.Neighborhood,
		Rua:    response.Street,
	}, nil
}

func (v *ViaCepAPI) BuscaByCEP(cep string) (Endereco, error) {
	resp, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return Endereco{}, fmt.Errorf("erro ao consumir api %v", err)
	}

	type cepResponse struct {
		Logradouro string `json:"logradouro"`
		Bairro     string `json:"bairro"`
		Localidade string `json:"localidade"`
		Estado     string `json:"estado"`
	}

	var response cepResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Endereco{}, fmt.Errorf("erro ao realizar o decode %v", err)
	}

	return Endereco{
		Estado: response.Estado,
		Cidade: response.Localidade,
		Bairro: response.Bairro,
		Rua:    response.Logradouro,
	}, nil
}
