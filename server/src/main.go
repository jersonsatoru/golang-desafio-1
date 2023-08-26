package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
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

type ResponseServiceURL struct {
	USDBRL *USDBRL `json:"USDBRL"`
}

const DB = "cotacoes.db"
const PORT = ":8080"

func init() {
	conn, err := sql.Open("sqlite3", DB)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ctxSQL, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	sql := `
		CREATE TABLE IF NOT EXISTS cotacoes  (
			code VARCHAR,
			codein VARCHAR,
			name VARCHAR,
			high VARCHAR,
			low VARCHAR,
			varBid VARCHAR,
			pctChange VARCHAR,
			bid VARCHAR,
			ask VARCHAR,
			timestamp VARCHAR,
			create_date VARCHAR
		);
	`
	_, err = conn.ExecContext(ctxSQL, sql, "")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		serviceURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Millisecond*200))
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceURL, nil)
		client := http.Client{}
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()
		data, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		var responseServiceURL ResponseServiceURL
		err = json.Unmarshal(data, &responseServiceURL)
		if err != nil {
			panic(err)
		}

		conn, err := sql.Open("sqlite3", DB)
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		ctxSQL, cancel := context.WithTimeout(r.Context(), time.Duration(time.Nanosecond*100))
		defer cancel()

		sql := `
			INSERT INTO cotacoes
			(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
			VALUES (?, ?,?,?,?,?,?,?,?,?,?)
		`
		_, err = conn.ExecContext(ctxSQL, sql,
			responseServiceURL.USDBRL.Code, responseServiceURL.USDBRL.Codein, responseServiceURL.USDBRL.Name,
			responseServiceURL.USDBRL.High, responseServiceURL.USDBRL.Low, responseServiceURL.USDBRL.VarBid,
			responseServiceURL.USDBRL.PctChange, responseServiceURL.USDBRL.Bid, responseServiceURL.USDBRL.Ask,
			responseServiceURL.USDBRL.Timestamp, responseServiceURL.USDBRL.CreateDate,
		)
		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{"bid": responseServiceURL.USDBRL.Bid})
		if err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server started at %s", PORT)

}
