package storages

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"syscall"
	"time"
)

// NewStorage returns a file system storage engine
func NewFileSystemStorage(location string) Storage {
	return &FileSystemStorage{
		Location: location,
	}
}

func (s *FileSystemStorage) NewFromParams(params map[string]string) (Storage, error) {
	return NewFileSystemStorage(params["location"]), nil
}

// Storage is a file system storage handler
type FileSystemStorage struct {
	Location string
}

// Save saves a file at the given path
func (s *FileSystemStorage) Save(filepath string, content []byte) error {
	return s.SaveWithPermissions(filepath, content, DefaultFilePermissions)
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

// AccessedTime returns the last access time.
func (s *FileSystemStorage) AccessedTime(filepath string) (time.Time, error) {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return time.Time{}, err
	}
	if runtime.GOOS == "windows" {
		return s.ModifiedTime(filepath)
	}
	stat := fi.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec)), nil
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
	return time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)), nil
}

// Size returns the size of the given file
func (s *FileSystemStorage) Size(filepath string) int64 {
	fi, err := os.Stat(s.Path(filepath))
	if err != nil {
		return 0
	}
	return fi.Size()
}
