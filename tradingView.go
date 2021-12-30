package main

import (
	"encoding/json"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"io"
	"log"
	"strings"
)

type TradingViewMessage struct {
	Ticker   string
	Price    float64
	Volume   string
	Exchange string
	Action   alpaca.Side
	Interval string
}

func unmarshalTradingViewBody(body *io.ReadCloser) (*TradingViewMessage, error) {
	var t *TradingViewMessage
	// Decode Request Body into TradingViewMessage
	err := json.NewDecoder(*body).Decode(&t)
	if err != nil {
		return nil, err
	}

	t.Ticker = sanitizeInput(t.Ticker)
	t.Exchange = sanitizeInput(t.Exchange)
	t.Interval = sanitizeInput(t.Interval)
	t.Volume = sanitizeInput(t.Volume)

	log.Println("Ticker:", t.Ticker)
	log.Println("Price:", t.Price)
	log.Println("Action:", t.Action)
	return t, err
}

func sanitizeInput(s string) string {
	l := strings.Replace(s, "\n", "", -1)
	l = strings.Replace(l, "\r", "", -1)
	return l
}
