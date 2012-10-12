package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

var targetServer = "www.taobao.com:80"

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://" + targetServer + r.URL.String())
	if err != nil {
		return
	}
	io.Copy(w, resp.Body)
	r.Body.Close()
}

func main() {
	log.Printf("Starting http Server ... ")
	http.HandleFunc("/", handler)
	s := &http.Server{
		Addr:           ":12345",
		Handler:        nil,
		ReadTimeout:    100 * time.Millisecond,
		WriteTimeout:   1000 * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
