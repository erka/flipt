package object

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
	"go.flipt.io/flipt/internal/containers"
	"go.flipt.io/flipt/internal/storage"
	"go.flipt.io/flipt/internal/storage/fs"
	"go.uber.org/zap/zaptest"
)

const testS3Bucket = "testdata"

var minioURL = os.Getenv("TEST_S3_ENDPOINT")

func Test_S3_Store(t *testing.T) {
	ch := make(chan struct{})
	store, skip := testStore(t, WithPollOptions(
		fs.WithInterval(time.Second),
		fs.WithNotify(t, func(modified bool) {
			if modified {
				close(ch)
			}
		}),
	))
	if skip {
		return
	}

	require.Equal(t, "s3", store.String())

	// flag shouldn't be present until we update it
	require.Error(t, store.View(func(s storage.ReadOnlyStore) error {
		_, err := s.GetFlag(context.TODO(), "production", "foo")
		return err
	}), "flag should not be defined yet")

	updated := []byte(`namespace: production
flags:
    - key: foo
      name: Foo`)

	buf := bytes.NewReader(updated)

	s3Client := testS3Client(t, "minio", minioURL)
	// update features.yml
	path := "features.yml"
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &store.bucket,
		Key:    &path,
		Body:   buf,
	})
	require.NoError(t, err)

	// assert matching state
	select {
	case <-ch:
	case <-time.After(time.Minute):
		t.Fatal("timed out waiting for update")
	}

	t.Log("received new snapshot")

	require.NoError(t, store.View(func(s storage.ReadOnlyStore) error {
		_, err = s.GetNamespace(context.TODO(), "production")
		if err != nil {
			return err
		}

		_, err = s.GetFlag(context.TODO(), "production", "foo")
		if err != nil {
			return err
		}

		_, err = s.GetNamespace(context.TODO(), "prefix")
		return err
	}))

}

func Test_S3_Store_WithPrefix(t *testing.T) {
	store, skip := testStore(t, WithBucket(testS3Bucket, "prefix"))
	if skip {
		return
	}

	// namespace shouldn't exist as it has been filtered out by the prefix
	require.Error(t, store.View(func(s storage.ReadOnlyStore) error {
		_, err := s.GetNamespace(context.TODO(), "production")
		return err
	}), "production namespace shouldn't be retrieavable")
}

func testStore(t *testing.T, opts ...containers.Option[SnapshotStore]) (*SnapshotStore, bool) {
	t.Helper()

	if minioURL == "" {
		t.Skip("Set non-empty TEST_S3_ENDPOINT env var to run this test.")
		return nil, true
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	source, err := NewSnapshotStore(ctx, zaptest.NewLogger(t),
		append([]containers.Option[SnapshotStore]{
			WithFS("s3", S3FSURL(testS3Bucket, "minio", minioURL)),
			WithBucket(testS3Bucket, ""),
		},
			opts...)...,
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = source.Close()
	})

	return source, false
}

func testS3Client(t *testing.T, region string, endpoint string) *s3.Client {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region))
	require.NoError(t, err)
	var s3Opts []func(*s3.Options)
	if endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = &endpoint
			o.UsePathStyle = true
			o.Region = region
		})
	}
	return s3.NewFromConfig(cfg, s3Opts...)
}
