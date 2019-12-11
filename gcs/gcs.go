package gcs

import (
	"context"
	"io"
	"mime"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/ulule/gostorages"
	"google.golang.org/api/option"
)

// Storage is a gcs storage.
type Storage struct {
	bucket *storage.BucketHandle
}

// NewStorage returns a new Storage.
func NewStorage(ctx context.Context, credentialsFile, bucket string) (*Storage, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &Storage{bucket: client.Bucket(bucket)}, nil
}

// Save saves content to path.
func (g *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	w := g.bucket.Object(path).NewWriter(ctx)
	w.ContentType = mime.TypeByExtension(filepath.Ext(path))

	if _, err := io.Copy(w, content); err != nil {
		return err
	}
	return w.Close()
}

// Stat returns path metadata.
func (g *Storage) Stat(ctx context.Context, path string) (*gostorages.Stat, error) {
	attrs, err := g.bucket.Object(path).Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, gostorages.ErrNotExist
	} else if err != nil {
		return nil, err
	}

	return &gostorages.Stat{
		ModifiedTime: attrs.Updated,
		Size:         attrs.Size,
	}, nil
}

// Open opens path for reading.
func (g *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return g.bucket.Object(path).NewReader(ctx)
}

// Delete deletes path.
func (g *Storage) Delete(ctx context.Context, path string) error {
	return g.bucket.Object(path).Delete(ctx)
}
