package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Cotacao struct {
	ValorCotacao string `json:"bid"`
}

const (
	cotacaoURL = "http://localhost:8080/cotacao"
	outputFile = "cotacao.txt"
)

func main() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", cotacaoURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar nova requisição: %v\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao realizar a requisição: %v\n", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao realizar o parse da resposta: %v\n", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("Dólar: %s", cotacao.ValorCotacao))

}
