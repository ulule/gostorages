package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/ulule/gostorages"
	"google.golang.org/api/option"
	"io"
	"mime"
	"path/filepath"
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
func (g *Storage) Save(ctx context.Context, content io.Reader, path string) (rerr error) {
	w := g.bucket.Object(path).NewWriter(ctx)
	w.ContentType = mime.TypeByExtension(filepath.Ext(path))

	defer func() {
		if err := w.Close(); err != nil {
			rerr = err
		}
	}()

	if _, err := io.Copy(w, content); err != nil {
		return err
	}

	return rerr
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
	r, err := g.bucket.Object(path).NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, gostorages.ErrNotExist
	}

	return r, err
}

// Delete deletes path.
func (g *Storage) Delete(ctx context.Context, path string) error {
	return g.bucket.Object(path).Delete(ctx)
}

// OpenWithStat opens path for reading with file stats.
func (g *Storage) OpenWithStat(ctx context.Context, path string) (io.ReadCloser, *gostorages.Stat, error) {
	r, err := g.bucket.Object(path).NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, nil, gostorages.ErrNotExist
	}
	return r, &gostorages.Stat{
		ModifiedTime: r.Attrs.LastModified,
		Size:         r.Attrs.Size,
	}, err
}
