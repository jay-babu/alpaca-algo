package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/shopspring/decimal"
	"log"
)

const (
	ReserveBalance = 25000
)

func getAccountBalance() decimal.Decimal {
	// Get our account information.
	account, err := alpaca.GetAccount()

	if err != nil {
		panic(err)
	}

	// Check if our account is restricted from trading.
	if account.TradingBlocked {
		log.Fatalln("Account is currently restricted from trading.")
	}

	// Check how much money we can use to open new positions.
	log.Printf("$%v is available as buying power.\n", account.Cash)

	return account.Cash
}

func canBuy(ticker string) bool {
	quote, err := marketdata.GetLatestQuote(ticker)
	if err != nil {
		log.Println("Can't get latest price for", ticker+":", err)
		return false
	}
	balance := getAccountBalance()
	// Subtract ReserveBalance and Ticker BidPrice + Buffer
	balance = balance.Sub(decimal.NewFromFloat(ReserveBalance)).Sub(decimal.NewFromFloat(quote.BidPrice * 1.01))
	return balance.GreaterThan(decimal.Zero)
}
