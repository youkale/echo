package main

import (
	"encoding/json"
	"io"
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

func main() {
	http.ListenAndServe(":6111", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			w.Write(res)
		}
	}))
}
