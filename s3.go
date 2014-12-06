package storages

import (
	"errors"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"mime"
	"net/http"
	"path"
)

var ACLs = map[string]s3.ACL{
	"private":                   s3.Private,
	"public-read":               s3.PublicRead,
	"public-read-write":         s3.PublicReadWrite,
	"authenticated-read":        s3.AuthenticatedRead,
	"bucket-owner-read":         s3.BucketOwnerRead,
	"bucket-owner-full-control": s3.BucketOwnerFull,
}

func NewS3Storage(accessKeyId string, secretAccessKey string, bucketName string, location string, region aws.Region, acl s3.ACL) (Storage, error) {
	return &S3Storage{
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
		Location:        location,
		Region:          region,
		ACL:             acl,
	}, nil
}

func (s *S3Storage) Params(params map[string]string) error {
	ACL, ok := ACLs[params["acl"]]

	if !ok {
		return errors.New(fmt.Sprintf("The ACL %s does not exist", params["acl"]))
	}

	Region, ok := aws.Regions[params["region"]]

	if !ok {
		return errors.New(fmt.Sprintf("The Region %s does not exist", params["region"]))
	}

	s.AccessKeyId = params["access_key_id"]
	s.SecretAccessKey = params["secret_access_key"]
	s.BucketName = params["bucket_name"]
	s.Location = params["location"]
	s.Region = Region
	s.ACL = ACL

	return nil
}

type S3Storage struct {
	AccessKeyId     string
	SecretAccessKey string
	BucketName      string
	Location        string
	Region          aws.Region
	ACL             s3.ACL
}

// Auth returns a Auth instance
func (s *S3Storage) Auth() (auth aws.Auth, err error) {
	return aws.GetAuth(s.AccessKeyId, s.SecretAccessKey)
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

// Exists checks if the given file is in the bucket
func (s *S3Storage) Exists(filepath string) bool {
	bucket, err := s.Bucket()

	if err != nil {
		return false
	}

	res, err := bucket.Head(s.Path(filepath))

	if err != nil {
		return false
	}

	if res.StatusCode != http.StatusOK {
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

func (s *S3Storage) Save(filepath string, content []byte) error {
	return s.SaveWithContentType(filepath, content, mime.TypeByExtension(filepath))
}

// Path joins the given file to the storage output directory
func (s *S3Storage) Path(filepath string) string {
	return path.Join(s.Location, filepath)
}
