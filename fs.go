package storages

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// NewStorage returns a file system storage engine
func NewFileSystemStorage(location string, baseURL string) Storage {
	return &FileSystemStorage{
		Location: location,
		BaseURL:  baseURL,
	}
}

// Storage is a file system storage handler
type FileSystemStorage struct {
	Location string
	BaseURL  string
}

// Save saves a file at the given path
func (s *FileSystemStorage) Save(filepath string, content []byte) error {
	return s.SaveWithPermissions(filepath, content, DefaultFilePermissions)
}

func (s *FileSystemStorage) URL(filename string) string {
	if s.HasBaseURL() {
		return strings.Join([]string{s.BaseURL, s.Path(filename)}, "/")
	}

	return ""
}

func (s *FileSystemStorage) HasBaseURL() bool {
	return s.BaseURL != ""
}

// SaveWithPermissions saves a file with the given permissions to the storage
func (s *FileSystemStorage) SaveWithPermissions(filepath string, content []byte, perm os.FileMode) error {
	_, err := os.Stat(s.Location)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.Path(filepath), content, perm)

	return err
}

// Open returns the file content
func (s *FileSystemStorage) Open(filepath string) ([]byte, error) {
	return ioutil.ReadFile(s.Path(filepath))
}

// Delete the file from storage
func (s *FileSystemStorage) Delete(filepath string) error {
	return os.Remove(s.Path(filepath))
}

// ModifiedTime returns the last update time
func (s *FileSystemStorage) ModifiedTime(filepath string) (time.Time, error) {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}

	return fi.ModTime(), nil
}

// Path joins the given file to the storage path
func (s *FileSystemStorage) Path(filepath string) string {
	return path.Join(s.Location, filepath)
}

// Exists checks if the given file is in the storage
func (s *FileSystemStorage) Exists(filepath string) bool {
	_, err := os.Stat(s.Path(filepath))
	return err == nil
}

func (s *FileSystemStorage) AccessedTime(filepath string) (time.Time, error) {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}
	if runtime.GOOS == "windows" {
		return s.ModifiedTime(filepath)
	}

	stat := fi.Sys().(*syscall.Stat_t)

	atim := reflect.ValueOf(stat).FieldByName("Atim")

	sec := atim.FieldByName("Sec").Int()

	nsec := atim.FieldByName("Nsec").Int()

	return time.Unix(int64(sec), int64(nsec)), nil
}

// CreatedTime returns the last access time.
func (s *FileSystemStorage) CreatedTime(filepath string) (time.Time, error) {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}
	if runtime.GOOS == "windows" {
		return s.ModifiedTime(filepath)
	}
	stat := fi.Sys().(*syscall.Stat_t)

	atim := reflect.ValueOf(stat).FieldByName("Ctim")

	sec := atim.FieldByName("Sec").Int()

	nsec := atim.FieldByName("Nsec").Int()

	return time.Unix(int64(sec), int64(nsec)), nil
}

// Size returns the size of the given file
func (s *FileSystemStorage) Size(filepath string) int64 {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return 0
	}
	return fi.Size()
}
