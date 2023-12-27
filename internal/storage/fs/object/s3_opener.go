package object

import (
	"context"
	"fmt"
	"net/url"

	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	gcaws "gocloud.dev/aws"
	gcblob "gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

const (
	S3Schema = "s3i"
)

func init() {
	gcblob.DefaultURLMux().RegisterBucket(S3Schema, new(urlSessionOpener))
}

type urlSessionOpener struct{}

func (o *urlSessionOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*gcblob.Bucket, error) {
	cfg, err := gcaws.V2ConfigFromURLParams(ctx, u.Query())
	if err != nil {
		return nil, fmt.Errorf("open bucket %v: %w", u, err)
	}
	clientV2 := s3v2.NewFromConfig(cfg, func(o *s3v2.Options) {
		o.UsePathStyle = true
	})
	return s3blob.OpenBucketV2(ctx, clientV2, u.Host, &s3blob.Options{})
}
