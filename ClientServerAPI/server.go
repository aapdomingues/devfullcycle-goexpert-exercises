package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

const (
	cotacaoURL  = "https://economia.awesomeapi.com.br/json/last"
	moedaUSDBRL = "USD-BRL"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", CotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {

	c, err := getCotacao(moedaUSDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Cotacao:", c)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func getCotacao(moeda string) (*Cotacao, error) {
	req, err := http.Get(cotacaoURL + "/" + moeda)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var mc map[string]Cotacao
	err = json.Unmarshal(res, &mc)
	if err != nil {
		return nil, err
	}

	fmt.Println("mapa: ", mc)

	c := mc["USDBRL"]

	return &c, nil

}
