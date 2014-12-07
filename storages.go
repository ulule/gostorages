package storages

import (
	"time"
)

type Storage interface {
	Save(filepath string, content []byte) error
	Path(filepath string) string
	Exists(filepath string) bool
	Delete(filepath string) error
	Open(filepath string) ([]byte, error)
	ModifiedTime(filepath string) (time.Time, error)
	Size(filepath string) int64
	URL(filename string) string
	HasBaseURL() bool
}
