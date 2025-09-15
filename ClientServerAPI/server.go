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
	cotacaoURL  = "https://economia.awesomeapi.com.br/json/last"
	moedaUSDBRL = "USD-BRL"
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

	// db, err := sql.Open("sqlite3", "./app.db")
	// if err != nil {
	// 	panic(err)
	// }

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
	mux.HandleFunc("/", handler.CotacaoHandlerFunc)
	http.ListenAndServe(":8080", mux)
}

func (h *Handler) CotacaoHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	c, err := getCotacao(ctx, moedaUSDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Cotacao:", c)

	ctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	err = h.service.SalvarCotacao(ctx, c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func getCotacao(ctx context.Context, moeda string) (*Cotacao, error) {
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

func (s *Service) SalvarCotacao(ctx context.Context, cotacao *Cotacao) error {
	return s.repository.SalvarCotacao(ctx, cotacao)
}

func (r *Repository) SalvarCotacao(ctx context.Context, cotacao *Cotacao) error {
	return gorm.G[Cotacao](r.database).Create(ctx, cotacao)
}
