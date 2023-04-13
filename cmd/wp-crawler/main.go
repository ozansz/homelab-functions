package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/ozansz/homelab-functions/pkg/minioext"
	"github.com/ozansz/homelab-functions/pkg/wordpress"
)

const (
	minioAccessKeyIDEnv     = "MINIO_ACCESS_KEY_ID"
	minioSecretAccessKeyEnv = "MINIO_SECRET_ACCESS_KEY"
)

var (
	url         = flag.String("url", "", "URL of the WordPress site to crawl")
	timeout     = flag.Duration("timeout", 5*time.Minute, "Timeout for crawler")
	httpTimeout = flag.Duration("http-timeout", 10*time.Second, "Timeout for HTTP requests")

	minioEndpoint    = flag.String("minio-endpoint", "", "Minio endpoint")
	minioRegion      = flag.String("minio-region", "", "Minio region")
	minioBucket      = flag.String("minio-bucket", "", "Minio bucket")
	minioHTTPTimeout = flag.Duration("minio-http-timeout", 10*time.Second, "Timeout for Minio HTTP requests")

	debugOutput = flag.Bool("debug-output", false, "Debug output")

	minioAccessKeyID     string
	minioSecretAccessKey string
	minioCl              *minioext.Client
)

func main() {
	flag.Parse()
	minioAccessKeyID = os.Getenv(minioAccessKeyIDEnv)
	minioSecretAccessKey = os.Getenv(minioSecretAccessKeyEnv)

	mustValidateConfig()

	ctx := context.Background()

	wpCl := wordpress.NewClient(*url, wordpress.WithTimeout(*httpTimeout))

	var err error
	if !*debugOutput {
		minioCl, err = minioext.NewClient(*minioEndpoint, *minioRegion, minioext.WithTimeout(*minioHTTPTimeout), minioext.WithCredentials(minioAccessKeyID, minioSecretAccessKey))
		if err != nil {
			log.Fatalf("failed to create minio client: %v", err)
		}
	}

	wpData, err := wpCl.GetAll(ctx)
	if err != nil {
		log.Fatalf("failed to get all data from WordPress: %v", err)
	}

	data, err := wpData.Marshal()
	if err != nil {
		log.Fatalf("failed to marshal data: %v", err)
	}

	if *debugOutput {
		log.Printf("data: %s", data)
		return
	}

	if err := minioCl.BatchUploadBytesWithDateTimePath(ctx, *minioBucket, data, minioext.LayoutYYYYMMDDHHMM, minio.PutObjectOptions{
		ContentType: "application/json",
	}); err != nil {
		log.Fatalf("failed to upload data to minio: %v", err)
	}

	log.Println("ok!")
}

func mustValidateConfig() {
	if *url == "" {
		log.Fatal("url is required")
	}
	if !*debugOutput {
		if *minioEndpoint == "" {
			log.Fatal("minio-endpoint is required")
		}
		if *minioRegion == "" {
			log.Fatal("minio-region is required")
		}
		if *minioBucket == "" {
			log.Fatal("minio-bucket is required")
		}
		if minioAccessKeyID == "" {
			log.Fatalf("%s is required", minioAccessKeyIDEnv)
		}
		if minioSecretAccessKey == "" {
			log.Fatalf("%s is required", minioSecretAccessKeyEnv)
		}
	}
}
