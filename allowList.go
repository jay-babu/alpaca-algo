package main

import "errors"

var AllowSet = map[string]struct{}{
	"NVDA": {},
	"AMD":  {},
}

func allowedTicker(ticker string) error {
	_, ok := AllowSet[ticker]

	var err error

	if !ok {
		err = errors.New(ticker + "not allowed to be purchased")
		return err
	}

	return err
}
