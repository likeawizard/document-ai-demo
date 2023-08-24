package store

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/likeawizard/document-ai-demo/config"
	"google.golang.org/api/option"
)

type GCloudBucket struct {
	bucket string
}

func NewGCloudStore(cfg config.StorageCfg) *GCloudBucket {
	return &GCloudBucket{
		bucket: cfg.Location,
	}
}

func (gcStore *GCloudBucket) Get(filename string) (io.ReadCloser, error) {
	ctx := context.Background()
	bkt, err := gcStore.getBucket(ctx)
	if err != nil {
		return nil, err
	}

	obj := bkt.Object(filename)
	return obj.NewReader(ctx)
}

func (gcStore *GCloudBucket) Store(filename string, r io.Reader) error {
	ctx := context.Background()
	bkt, err := gcStore.getBucket(ctx)
	if err != nil {
		return err
	}

	obj := bkt.Object(filename)
	w := obj.NewWriter(ctx)
	br := bufio.NewReader(r)
	_, err = br.WriteTo(w)
	if err != nil {
		return err
	}

	return w.Close()
}

// TODO: this relies on bucket/objects being public. Could generate a temporary SignedURL for more a more robust solution. All files currently are publicly available which is a big no-no for real data.
func (gcStore *GCloudBucket) GetURL(filename string) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcStore.bucket, filename), nil

}

func (gcStore *GCloudBucket) getBucket(ctx context.Context) (*storage.BucketHandle, error) {
	auth := option.WithCredentialsFile("document-ai-creds.json")
	client, err := storage.NewClient(ctx, auth)
	if err != nil {
		return nil, err
	}
	return client.Bucket(gcStore.bucket), nil
}
