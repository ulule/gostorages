package fs

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ulule/gostorages"
)

func Test(t *testing.T) {
	dir, err := ioutil.TempDir("", "gostorages-")
	if err != nil {
		t.Fatal(err)
	}
	storage := NewStorage(Config{Root: dir})
	ctx := context.Background()

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
