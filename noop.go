package gostorages

import (
	"context"
	"io"
)

// NewNoop returns a new no-op storage.
func NewNoop() Storage {
	return noop{}
}

type noop struct{}

func (noop) Save(ctx context.Context, content io.Reader, path string) error { return nil }
func (noop) Stat(ctx context.Context, path string) (*Stat, error)           { return nil, nil }
func (noop) Open(ctx context.Context, path string) (io.ReadCloser, error)   { return nil, nil }
func (noop) OpenWithStat(ctx context.Context, path string) (io.ReadCloser, *Stat, error) {
	return nil, nil, nil
}
func (noop) Delete(ctx context.Context, path string) error { return nil }
