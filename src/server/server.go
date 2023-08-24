package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type USDBRL struct {
	Code       string    `json:"code"`
	Codein     string    `json:"codein"`
	Name       string    `json:"name"`
	High       string    `json:"high"`
	Low        string    `json:"low"`
	VarBid     string    `json:"varBid"`
	PctChange  string    `json:"pctChange"`
	Bid        string    `json:"bid"`
	Ask        string    `json:"ask"`
	Timestamp  string    `json:"timestamp"`
	CreateDate time.Time `json:"create_date"`
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		serviceURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
		defer cancel()

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceURL, nil)
		if err != nil {
			panic(err)
		}

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

		var usdBrl *USDBRL
		err = json.Unmarshal(data, usdBrl)
		if err != nil {
			panic(err)
		}
	})

	http.ListenAndServe(":8080", mux)
}
