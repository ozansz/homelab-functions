package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
)

type Response struct {
	RemoteAddr string      `json:"remote_addr"`
	Method     string      `json:"method"`
	URI        URI         `json:"uri"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
}

type URI struct {
	Proto       string     `json:"proto"`
	Host        string     `json:"host"`
	Path        string     `json:"path"`
	QueryParams url.Values `json:"query_params"`
}

func main() {
	flag.Parse()

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorResponse(w, fmt.Sprintf("failed to read request body: %v", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		uri := URI{
			Path:        r.URL.Path,
			Proto:       r.Proto,
			QueryParams: r.URL.Query(),
			Host:        r.Host,
		}
		resp := Response{
			URI:        uri,
			Body:       string(body),
			Headers:    r.Header,
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			errorResponse(w, fmt.Sprintf("failed to encode response: %v", err))
			return
		}
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}

func errorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, message))); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}
