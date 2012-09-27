package main

import (
	"io"
	"log"
	"net/http"
)

type Proxy struct {
}

func NewProxy() *Proxy { return &Proxy{} }

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	switch r.Method {
	default:
		{
			log.Print("Cannot handle method ", r.Method)
			http.Error(wr, "501 I only handle GET and POST", http.StatusNotImplemented)
			return
		}
	case "GET":
		{
			log.Printf("getting %v", r.URL)
			resp, err = http.Get(r.URL.String())
		}
	case "POST":
		{
			resp, err = http.Post(r.URL.String(), r.Header["Content-Type"][0], r.Body)
			r.Body.Close()
		}
	}

	// combined for GET/POST
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", resp.Header["Content-Type"][0])
	wr.WriteHeader(resp.StatusCode)

	io.Copy(wr, resp.Body)

	resp.Body.Close()
}

func main() {
	proxy := NewProxy()
	err := http.ListenAndServe(":12345", proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
