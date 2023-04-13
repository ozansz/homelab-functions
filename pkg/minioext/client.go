package minioext

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	cl     *minio.Client
	region string
}

type NewClientOption func(*minio.Options)

func WithSSL() NewClientOption {
	return func(opts *minio.Options) {
		opts.Secure = true
	}
}

func WithTimeout(timeout time.Duration) NewClientOption {
	return func(opts *minio.Options) {
		opts.Transport.(*http.Transport).DialContext = (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
			DualStack: true,
		}).DialContext
	}
}

func WithHTTPTransport(transport *http.Transport) NewClientOption {
	return func(opts *minio.Options) {
		opts.Transport = transport
	}
}

func WithCredentials(accessKeyID, secretAccessKey string) NewClientOption {
	return func(opts *minio.Options) {
		opts.Creds = credentials.NewStaticV4(accessKeyID, secretAccessKey, "")
	}
}

func NewClient(endpoint, region string, opts ...NewClientOption) (*Client, error) {
	minioOpts := &minio.Options{
		Transport: defaultHTTPTransport(),
	}
	for _, opt := range opts {
		opt(minioOpts)
	}
	cl, err := minio.New(endpoint, minioOpts)
	if err != nil {
		return nil, err
	}
	return &Client{
		cl:     cl,
		region: region,
	}, nil
}

// Returns true if the bucket was created, false if it already exists.
// Returns an error if the bucket could not be created.
func (cl *Client) CreateBucketIfNotExists(ctx context.Context, bucket string) (bool, error) {
	exists, err := cl.cl.BucketExists(ctx, bucket)
	if exists {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := cl.CreateBucket(ctx, bucket); err != nil {
		return false, err
	}
	return true, nil
}

func (cl *Client) CreateBucket(ctx context.Context, bucket string) error {
	if err := cl.cl.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: cl.region}); err != nil {
		return err
	}
	return nil
}

func (cl *Client) UploadBytes(ctx context.Context, bucket, objectName string, data []byte, opts minio.PutObjectOptions) error {
	if _, err := cl.CreateBucketIfNotExists(ctx, bucket); err != nil {
		return err
	}
	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		return err
	}
	_, err = cl.cl.FPutObject(ctx, bucket, objectName, f.Name(), opts)
	return err
}

func (cl *Client) UploadBytesWithDatePath(ctx context.Context, bucket, objectName string, data []byte, opts minio.PutObjectOptions) error {
	now := time.Now().UTC()
	path := fmt.Sprintf("%s/%s", now.Format("2006/01/02"), objectName)

	return cl.UploadBytes(ctx, bucket, path, data, opts)
}

func (cl *Client) BatchUploadBytesWithDatePath(ctx context.Context, bucket string, objects map[string][]byte, opts minio.PutObjectOptions) error {
	now := time.Now().UTC()
	for objectName, data := range objects {
		path := fmt.Sprintf("%s/%s", now.Format("2006/01/02"), objectName)
		if err := cl.UploadBytes(ctx, bucket, path, data, opts); err != nil {
			return fmt.Errorf("failed to upload %s: %w", objectName, err)
		}
	}
	return nil
}

func defaultHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
