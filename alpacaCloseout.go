package main

import (
	"encoding/json"
	"errors"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"log"
	"sync"
	"time"
)

const (
	BufferTime = 15
)

func CloseoutPositionsBeforeMarketCloses() {
	clock, err := alpaca.GetClock()
	if err != nil {
		return
	}
	durationUntilClose := clock.NextClose.Add(time.Duration(-BufferTime) * time.Minute).Sub(time.Now())
	log.Println("Duration until Close", durationUntilClose)
	timer := time.NewTimer(durationUntilClose)
	for {
		select {
		case <-timer.C:
			clock, err := alpaca.GetClock()
			if err != nil {
				timer.Stop()
				return
			}
			log.Println("Duration until Close ", durationUntilClose)
			var wg sync.WaitGroup
			for ticker := range AllowSet {
				ticker := ticker
				wg.Add(1)
				go func() {
					_ = alpaca.DefaultClient.ClosePosition(ticker)
					defer wg.Done()
				}()
			}
			wg.Wait()

			durationUntilClose := clock.NextClose.Add(time.Duration(-BufferTime) * time.Minute).Sub(time.Now())
			timer.Reset(durationUntilClose)
		}
	}
}

func checkHours(allowed chan error) {
	clock, err := alpaca.GetClock()
	if err != nil {
		allowed <- err
		return
	}

	// No purchases / sales within 15 minutes of closing time or when closed
	if !clock.IsOpen {
		out, _ := json.Marshal(clock)
		err = errors.New("market is closed " + string(out))
		allowed <- err
		return
	}

	allowed <- err
}
