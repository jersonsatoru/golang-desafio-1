package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

type ResponseServer struct {
	Bid string `json:"bid"`
}

func main() {
	serviceURL := "http://localhost:8080/cotacao"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*300))
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

	var responseServer ResponseServer
	err = json.Unmarshal(data, &responseServer)
	if err != nil {
		panic(err)
	}

	line := fmt.Sprintf("Dolar: %s\n", responseServer.Bid)
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	file.WriteString(line)
}
