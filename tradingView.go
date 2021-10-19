package main

import (
	"encoding/json"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"io"
	"log"
)

type TradingViewMessage struct {
	Ticker   string
	Price    float64
	Volume   string
	Exchange string
	Action   alpaca.Side
	Interval string
}

func demarshalTradingViewBody(body *io.ReadCloser) (*TradingViewMessage, error) {
	var t *TradingViewMessage
	var err error
	// Decode Request Body into TradingViewMessage
	err = json.NewDecoder(*body).Decode(&t)
	if err != nil {
		return nil, err
	}

	log.Println("Ticker:", t.Ticker)
	log.Println("Price:", t.Price)
	log.Println("Action:", t.Action)
	return t, err
}
