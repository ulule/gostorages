package gcs

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/storage"
	gstorages "github.com/ulule/gostorages"
	"google.golang.org/api/option"
)

// TestOpenWithStatNonNotExistError ensures OpenWithStat returns the underlying
// error instead of panicking when NewReader fails with an error other than
// storage.ErrObjectNotExist (for example a canceled context).
func TestOpenWithStatNonNotExistError(t *testing.T) {
	client, err := storage.NewClient(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint("http://127.0.0.1:0"))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	g := &Storage{bucket: client.Bucket("bucket")}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r, stat, err := g.OpenWithStat(ctx, "object")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if errors.Is(err, gstorages.ErrNotExist) {
		t.Fatalf("expected the underlying error, got ErrNotExist: %v", err)
	}
	if r != nil || stat != nil {
		t.Fatalf("expected nil reader and stat on error, got %v %v", r, stat)
	}
}
