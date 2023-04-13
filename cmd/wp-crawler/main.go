package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ozansz/homelab-functions/pkg/wordpress"
)

var (
	url         = flag.String("url", "", "URL of the WordPress site to crawl")
	timeout     = flag.Duration("timeout", 5*time.Minute, "Timeout for crawler")
	httpTimeout = flag.Duration("http-timeout", 10*time.Second, "Timeout for HTTP requests")
)

func main() {
	flag.Parse()
	mustValidateFlags()

	ctx := context.Background()

	cl := wordpress.NewClient(*url, wordpress.WithTimeout(*httpTimeout))
	wpData, err := cl.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func mustValidateFlags() {
	if *url == "" {
		log.Fatal("url is required")
	}
}
