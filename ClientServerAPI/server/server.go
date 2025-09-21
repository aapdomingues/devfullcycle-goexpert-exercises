package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type CotacaoResponse struct {
	ValorCotacao string `json:"bid"`
}

type Cotacao struct {
	gorm.Model
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
	cotacaoURL = "https://economia.awesomeapi.com.br/json/last"
	USDBRL     = "USD-BRL"
)

type Repository struct {
	database *gorm.DB
}

func newRepository(db *gorm.DB) Repository {
	return Repository{
		database: db,
	}
}

type Service struct {
	repository Repository
}

func newService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

type Handler struct {
	service Service
}

func newHanlder(service Service) Handler {
	return Handler{
		service: service,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(&Cotacao{})

	repo := newRepository(db)
	srv := newService(repo)
	handler := newHanlder(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handler.CotacaoHandlerFunc)
	http.ListenAndServe(":8080", mux)
}

func (h *Handler) CotacaoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctxApi, cancelApi := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelApi()
	c, err := getCotacao(ctxApi, USDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctxDb, cancelDb := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDb()

	err = h.service.SalvarCotacao(ctxDb, c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cotacaoResp := CotacaoResponse{
		ValorCotacao: c.Bid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacaoResp)
}

func getCotacao(ctx context.Context, moeda string) (*Cotacao, error) {
	type ConversaoCodeMap map[string]string
	var internalApiCodeMap = ConversaoCodeMap{USDBRL: "USDBRL"}

	req, err := http.NewRequestWithContext(ctx, "GET", cotacaoURL+"/"+moeda, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var mc map[string]Cotacao
	err = json.Unmarshal(body, &mc)
	if err != nil {
		return nil, err
	}

	c := mc[internalApiCodeMap[moeda]]
	return &c, nil

}

func (s *Service) SalvarCotacao(ctx context.Context, cotacao *Cotacao) error {
	return s.repository.SalvarCotacao(ctx, cotacao)
}

func (r *Repository) SalvarCotacao(ctx context.Context, cotacao *Cotacao) error {
	return gorm.G[Cotacao](r.database).Create(ctx, cotacao)
}
