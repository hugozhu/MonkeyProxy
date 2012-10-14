package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"linkio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var link *linkio.Link
var random *rand.Rand

var dead_rate = flag.Int("dead_rate", 0, "死请求百分比")
var timeout_rate = flag.Int("timeout_rate", 0, "超时请求百分比")
var timeout = flag.Int("timeout", 1000, "超时设置")
var bad_rate = flag.Int("bad_rate", 0, "坏请求百分比")
var target = flag.String("target", "", "目标服务地址 ＊必填＊")
var port = flag.Int("port", 80, "本地端口")

func handler(w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	url := "http://" + *target + r.URL.String()
	randInt := random.Intn(100)

	block := false
	bad := false

	if *dead_rate > 0 {
		if randInt <= *dead_rate {
			log.Printf("%s dead response", url)
			time.Sleep(24 * time.Hour)
			block = true
			return
		}
	}

	if !block && *timeout_rate > 0 {
		if randInt > *dead_rate && randInt <= *dead_rate+*timeout_rate {
			log.Printf("%s [timeout]", url)
			time.Sleep(time.Duration(*timeout) * time.Millisecond)
			block = true
			return
		}
	}

	if !block && *bad_rate > 0 {
		if randInt > *dead_rate+*timeout_rate && randInt <= *dead_rate+*timeout_rate+*bad_rate {
			log.Printf("%s [bad]", url)
			bad = true
		}
	}

	if !block && !bad {
		log.Printf("%s [normal]", url)
	}

	req, _ := http.NewRequest(r.Method, url, r.Body)

	//copy client req headers to target server
	for k, v := range r.Header {
		for _, v1 := range v {
			req.Header.Add(k, v1)
		}
	}

	req.Header.Set("Host", *target)

	client := &http.Client{}

	resp, err = client.Do(req)

	if err != nil {
		return
	}

	//copy server response headers to client
	for k, v := range resp.Header {
		for _, v1 := range v {
			w.Header().Add(k, v1)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if bad {
		bytes, _ := ioutil.ReadAll(resp.Body)
		io.WriteString(w, string(bytes[:len(bytes)/2])) //only return partial content
	} else {
		io.Copy(w, resp.Body)
	}
	r.Body.Close()
}

func init() {
	link = linkio.NewLink(10240) //kbps
	random = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}

func main() {
	flag.Parse()
	if *target == "" {
		flag.Usage()
		os.Exit(0)
	}
	log.Printf("Starting http proxy server \nlocal port:\t%d\ntarget:\t%s\ntimeout rate:\t%d%%\ntimeout:\t%d毫秒\ndead rate:\t%d%%\nbad rate:\t%d%%",
		*port, *target, *timeout_rate, *timeout, *dead_rate, *bad_rate)

	http.HandleFunc("/", handler)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        nil,
		ReadTimeout:    1000 * time.Millisecond,
		WriteTimeout:   10000 * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
