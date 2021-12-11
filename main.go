package main

import (
	"errors"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"os"
)

func getEnv() (string, string, string) {
	secret := os.Getenv("APCA_API_SECRET_KEY")
	apiKey := os.Getenv("APCA_API_KEY_ID")
	baseUrl := os.Getenv("APCA_API_BASE_URL")

	if len(secret) == 0 || len(apiKey) == 0 || len(baseUrl) == 0 {
		panic(errors.New("env var is undefined"))
	}

	return secret, apiKey, baseUrl
}

func handleBuyAndSell(w http.ResponseWriter, req *http.Request) {

	ok := setupResponse(&w, req)
	if !ok {
		return
	}
	if (*req).Method == "OPTIONS" {
		return
	}

	transactionAllow := make(chan error, 1)
	enoughFundsAvail := make(chan bool, 1)
	t, err := unmarshalTradingViewBody(&req.Body)
	if err != nil {
		dumpBody(req)
		log.Println("Body has an error!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	go dumpBody(req)

	go checkHours(transactionAllow)
	go func() {
		enoughFundsAvail <- canBuy(t.Ticker)
	}()

	if err := <-transactionAllow; err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Error: %s", err)
		return
	}

	if enoughFunds := <-enoughFundsAvail; !enoughFunds {
		log.Println("Can't afford to buy:", t.Ticker)
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, "Can't afford to buy: %s", t.Ticker)
		return
	}

	_, err = placeOrder(t.Ticker, t.Action)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "Error: %s", err)
		return
	}
}

func placeOrder(ticker string, side alpaca.Side) (*alpaca.Order, error) {
	log.Println("Placing Order")

	if err := allowedTicker(ticker); err != nil {
		return nil, err
	}

	var placeOrderRequest alpaca.PlaceOrderRequest

	// 1 Share of Ticker
	qty := decimal.NewFromInt(1)

	switch side {
	case alpaca.Buy, alpaca.Sell:
		// Amount of Ticker currently owned
		_ = alpaca.DefaultClient.ClosePosition(ticker)
		status := make(chan string, 1)
		go func() {
			status <- checkTickerFilled(ticker)
		}()
		<-status

		log.Println("Closed", ticker)
		log.Println(side, "quantity", qty, "of", ticker)
		placeOrderRequest = alpaca.PlaceOrderRequest{
			AssetKey:    &ticker,
			Qty:         &qty,
			Type:        alpaca.Market,
			TimeInForce: alpaca.Day,
			Side:        side,
		}
	default:
		return nil, errors.New("action not supported. must be buy or sell")
	}

	return alpaca.PlaceOrder(placeOrderRequest)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) bool {
	s := make(map[string]struct{})
	s["52.89.214.238"] = struct{}{}
	s["34.212.75.30"] = struct{}{}
	s["54.218.53.128"] = struct{}{}
	s["52.32.178.7"] = struct{}{}
	req.Header.Get("X-Forwarded-For")
	_, ok := s[req.Header.Get("X-Forwarded-For")]
	if ok {
		(*w).Header().Set("Access-Control-Allow-Origin", req.Header.Get("X-Forwarded-For"))
	}
	fmt.Print("HOST: ")
	fmt.Println(req.Host)
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Host, User-Agent")
	return ok
}

func main() {
	getEnv()
	getAccountBalance()

	http.HandleFunc("/", handleBuyAndSell)

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println(err)
	}
}
