package fs

import (
	"context"
	"github.com/pkg/errors"
	"github.com/ulule/gostorages"
	"io"
	"os"
	"path/filepath"
)

// Storage is a filesystem storage.
type Storage struct {
	root string
}

// NewStorage returns a new filesystem storage.
func NewStorage(cfg Config) *Storage {
	return &Storage{root: cfg.Root}
}

// Config is the configuration for Storage.
type Config struct {
	Root string
}

func (fs *Storage) abs(path string) string {
	return filepath.Join(fs.root, path)
}

// Save saves content to path.
func (fs *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	abs := fs.abs(path)
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return errors.WithStack(err)
	}

	w, err := os.Create(abs)
	if err != nil {
		return errors.WithStack(err)
	}
	defer w.Close()

	if _, err := io.Copy(w, content); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Stat returns path metadata.
func (fs *Storage) Stat(ctx context.Context, path string) (*gostorages.Stat, error) {
	fi, err := os.Stat(fs.abs(path))
	if os.IsNotExist(err) {
		return nil, gostorages.ErrNotExist
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return &gostorages.Stat{
		ModifiedTime: fi.ModTime(),
		Size:         fi.Size(),
	}, nil
}

// Open opens path for reading.
func (fs *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	f, err := os.Open(fs.abs(path))
	if os.IsNotExist(err) {
		return nil, gostorages.ErrNotExist
	}
	return f, errors.WithStack(err)
}

// Delete deletes path.
func (fs *Storage) Delete(ctx context.Context, path string) error {
	return os.Remove(fs.abs(path))
}

// OpenWithStat opens path for reading with file stats.
func (fs *Storage) OpenWithStat(ctx context.Context, path string) (io.ReadCloser, *gostorages.Stat, error) {
	f, err := os.Open(fs.abs(path))
	if os.IsNotExist(err) {
		return nil, nil, gostorages.ErrNotExist
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return f, &gostorages.Stat{
		ModifiedTime: stat.ModTime(),
		Size:         stat.Size(),
	}, nil
}
