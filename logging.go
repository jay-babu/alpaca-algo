package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func dumpBody(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(requestDump))
}
