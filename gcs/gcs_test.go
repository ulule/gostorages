package gcs

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
	credFile := os.Getenv("CRED_FILE")
	bucket := os.Getenv("GCP_BUCKET")

	if credFile == "" || bucket == "" {
		t.SkipNow()
	}
	ctx := context.Background()
	storage, err := NewStorage(ctx, credFile, bucket)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = storage.Stat(ctx, "doesnotexist"); !errors.Is(err, gostorages.ErrNotExist) {
		t.Errorf("expected not exists, got %v", err)
	}

	before := time.Now()
	if err := storage.Save(ctx, bytes.NewBufferString("hello"), "world"); err != nil {
		t.Fatal(err)
	}
	now := time.Now()

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
