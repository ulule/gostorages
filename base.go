package storages

import (
	"path"
	"strings"
)

type BaseStorage struct {
	BaseURL  string
	Location string
}

func NewBaseStorage(location string, baseURL string) *BaseStorage {
	return &BaseStorage{
		BaseURL:  strings.TrimSuffix(baseURL, "/"),
		Location: location,
	}
}

func (s *BaseStorage) URL(filename string) string {
	if s.HasBaseURL() {
		return strings.Join([]string{s.BaseURL, s.Path(filename)}, "/")
	}

	return ""
}

func (s *BaseStorage) HasBaseURL() bool {
	return s.BaseURL != ""
}

// Path joins the given file to the storage path
func (s *BaseStorage) Path(filepath string) string {
	return path.Join(s.Location, filepath)
}
