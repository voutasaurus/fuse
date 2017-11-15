package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

func main() {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal("error creating request:", err)
	}
	res, err := (&http.Client{Transport: rt(h)}).Do(r)
	if err != nil {
		log.Fatal("error processing request:", err)
	}
	if res.StatusCode != 200 {
		log.Fatal("bad status:", res.StatusCode)
	}
	io.Copy(os.Stdout, res.Body)
}

var h http.HandlerFunc = serve

func serve(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world\n")
}

type funcRoundTripper func(r *http.Request) (*http.Response, error)

func (f funcRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func rt(h http.Handler) http.RoundTripper {
	return funcRoundTripper(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Result(), nil
	})
}
