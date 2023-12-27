package object

import (
	"net/url"
	"os"
	"strconv"

	"go.flipt.io/flipt/internal/storage/fs/object/blob"
	"gocloud.dev/blob/azureblob"
	"gocloud.dev/blob/gcsblob"
)

func AzureFSURL(container, endpoint string) (string, error) {
	if endpoint != "" {
		url, err := url.Parse(endpoint)
		if err != nil {
			return "", err
		}
		os.Setenv("AZURE_STORAGE_PROTOCOL", url.Scheme)
		os.Setenv("AZURE_STORAGE_IS_LOCAL_EMULATOR", strconv.FormatBool(url.Scheme == "http"))
		os.Setenv("AZURE_STORAGE_DOMAIN", url.Host)
	}
	return blob.StrUrl(azureblob.Scheme, container), nil
}

func S3FSURL(bucket, region, endpoint string) string {
	q := url.Values{}
	q.Set("awssdk", "v2")
	if region != "" {
		q.Set("region", region)
	}
	if endpoint != "" {
		q.Set("endpoint", endpoint)
	}
	return blob.StrUrl(S3Schema, bucket+"?"+q.Encode())
}

func GSFSURL(bucket string) string {
	return blob.StrUrl(gcsblob.Scheme, bucket)
}
