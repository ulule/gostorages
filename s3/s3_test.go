package s3

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ulule/gostorages"
)

func Test(t *testing.T) {
	accessKeyID := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")
	endpoint := os.Getenv("S3_ENDPOINT")
	acl := os.Getenv("S3_ACL")

	if accessKeyID == "" ||
		secretAccessKey == "" ||
		region == "" ||
		endpoint == "" ||
		acl == "" ||
		bucket == "" {
		t.SkipNow()
	}
	storage, err := NewStorage(Config{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		Bucket:          bucket,
		Endpoint:        &endpoint,
		ACL:             acl,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if _, err = storage.Stat(ctx, "doesnotexist"); !errors.Is(err, gostorages.ErrNotExist) {
		t.Errorf("expected not exists, got %v", err)
	}

	before := time.Now()
	if err := storage.Save(ctx, bytes.NewBufferString("hello"), "world"); err != nil {
		t.Fatal(err)
	}
	now := time.Now().Add(time.Second)

	stat, err := storage.Stat(ctx, "world")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if stat.Size != 5 {
		t.Errorf("expected size to be %d, got %d", 5, stat.Size)
	}
	if stat.ModifiedTime.Before(before) {
		t.Errorf("expected modtime to be after %v, got %v", before, stat.ModifiedTime)
	}
	if stat.ModifiedTime.After(now) {
		t.Errorf("expected modtime to be before %v, got %v", now, stat.ModifiedTime)
	}
}
