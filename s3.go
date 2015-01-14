package storages

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"mime"
	"path"
	"strings"
	"time"
)

var ACLs = map[string]s3.ACL{
	"private":                   s3.Private,
	"public-read":               s3.PublicRead,
	"public-read-write":         s3.PublicReadWrite,
	"authenticated-read":        s3.AuthenticatedRead,
	"bucket-owner-read":         s3.BucketOwnerRead,
	"bucket-owner-full-control": s3.BucketOwnerFull,
}

const LastModifiedFormat = time.RFC1123

func NewS3Storage(accessKeyId string, secretAccessKey string, bucketName string, location string, region aws.Region, acl s3.ACL, baseURL string) Storage {

	return &S3Storage{
		NewBaseStorage(location, baseURL),
		accessKeyId,
		secretAccessKey,
		bucketName,
		region,
		acl,
	}
}

type S3Storage struct {
	*BaseStorage
	AccessKeyId     string
	SecretAccessKey string
	BucketName      string
	Region          aws.Region
	ACL             s3.ACL
}

// Auth returns a Auth instance
func (s *S3Storage) Auth() (auth aws.Auth, err error) {
	return aws.GetAuth(s.AccessKeyId, s.SecretAccessKey)
}

func (s *S3Storage) URL(filename string) string {
	if s.HasBaseURL() {
		return strings.Join([]string{s.BaseURL, s.Path(filename)}, "/")
	}

	return ""
}

// Client returns a S3 instance
func (s *S3Storage) Client() (*s3.S3, error) {
	auth, err := s.Auth()

	if err != nil {
		return nil, err
	}

	return s3.New(auth, s.Region), nil
}

// Bucket returns a bucket instance
func (s *S3Storage) Bucket() (*s3.Bucket, error) {
	client, err := s.Client()

	if err != nil {
		return nil, err
	}

	return client.Bucket(s.BucketName), nil
}

// Open returns the file content in a dedicated bucket
func (s *S3Storage) Open(filepath string) ([]byte, error) {
	bucket, err := s.Bucket()

	if err != nil {
		return nil, err
	}

	return bucket.Get(s.Path(filepath))
}

// Delete the file from the bucket
func (s *S3Storage) Delete(filepath string) error {
	bucket, err := s.Bucket()

	if err != nil {
		return err
	}

	return bucket.Del(s.Path(filepath))
}

func (s *S3Storage) GetKey(filepath string) (*s3.Key, error) {
	bucket, err := s.Bucket()

	if err != nil {
		return nil, err
	}

	key, err := bucket.GetKey(s.Path(filepath))

	if err != nil {
		return nil, err
	}

	return key, nil
}

// Exists checks if the given file is in the bucket
func (s *S3Storage) Exists(filepath string) bool {
	_, err := s.GetKey(filepath)

	if err != nil {
		return false
	}

	return true
}

// Save saves a file at the given path in the bucket
func (s *S3Storage) SaveWithContentType(filepath string, content []byte, contentType string) error {
	bucket, err := s.Bucket()

	if err != nil {
		return err
	}

	err = bucket.Put(s.Path(filepath), content, contentType, s.ACL)

	return err
}

// Save saves a file at the given path in the bucket
func (s *S3Storage) Save(filepath string, content []byte) error {
	return s.SaveWithContentType(filepath, content, mime.TypeByExtension(filepath))
}

// Path joins the given file to the storage output directory
func (s *S3Storage) Path(filepath string) string {
	return path.Join(s.Location, filepath)
}

// Size returns the size of the given file
func (s *S3Storage) Size(filepath string) int64 {
	key, err := s.GetKey(filepath)

	if err != nil {
		return 0
	}

	return key.Size
}

// ModifiedTime returns the last update time
func (s *S3Storage) ModifiedTime(filepath string) (time.Time, error) {
	key, err := s.GetKey(filepath)

	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(LastModifiedFormat, key.LastModified)
}
