package s3

import (
	"context"
	"io"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ulule/gostorages"
)

// Storage is a s3 storage.
type Storage struct {
	bucket   string
	ACL      string
	s3       *s3.S3
	uploader *s3manager.Uploader
}

// NewStorage returns a new Storage.
func NewStorage(cfg Config) (*Storage, error) {
	s, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
		}),
		Region:   aws.String(cfg.Region),
		Endpoint: cfg.Endpoint,
	})
	if err != nil {
		return nil, err
	}

	return &Storage{
		bucket:   cfg.Bucket,
        ACL:      cfg.ACL,
		s3:       s3.New(s),
		uploader: s3manager.NewUploader(s),
	}, nil
}

// Config is the configuration for Storage.
type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	ACL             string
	Endpoint        *string
}

// Save saves content to path.
func (s *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	input := &s3manager.UploadInput{
		ACL:         aws.String(s.ACL),
		Body:        content,
		Bucket:      aws.String(s.bucket),
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(path))),
		Key:         aws.String(path),
	}

	_, err := s.uploader.UploadWithContext(ctx, input)
	return err
}

// Stat returns path metadata.
func (s *Storage) Stat(ctx context.Context, path string) (*gostorages.Stat, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}
	out, err := s.s3.HeadObjectWithContext(ctx, input)

	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
		return nil, gostorages.ErrNotExist
	} else if err != nil {
		return nil, err
	}

	return &gostorages.Stat{
		ModifiedTime: *out.LastModified,
		Size:         *out.ContentLength,
	}, nil
}

// Open opens path for reading.
func (s *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}
	out, err := s.s3.GetObjectWithContext(ctx, input)
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
		return nil, gostorages.ErrNotExist
	} else if err != nil {
		return nil, err
	}
	return out.Body, nil
}

// Delete deletes path.
func (s *Storage) Delete(ctx context.Context, path string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}
	_, err := s.s3.DeleteObjectWithContext(ctx, input)
	return err
}
