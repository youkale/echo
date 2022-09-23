package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
)

type Request struct {
	URL    string              `json:"url"`
	Query  string              `json:"query"`
	Header map[string][]string `json:"header"`
	Method string              `json:"method"`
	Host   string              `json:"host"`
	Body   string              `json:"body"`
}

var logger = log.Default()

var addr string

func init() {
	flag.StringVar(&addr, "l", ":6111", "listen address")
}

func main() {
	flag.Parse()
	logger.Print("listen port:", addr)
	http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := Request{
			URL:    r.RequestURI,
			Host:   r.Host,
			Method: r.Method,
			Header: r.Header,
			Query:  r.URL.RawQuery,
		}
		if nil != r.Body {
			if body, re := io.ReadAll(r.Body); nil == re {
				req.Body = string(body)
			}
		}

		if res, err := json.Marshal(req); nil == err {
			logger.Print("rec: " + string(res))
			w.Write(res)
		}
	}))
}
