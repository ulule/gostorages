gostorages
==========

A unified interface to manipulate storage backends in Go.

gostorages is used in [picfit](https://github.com/thoas/picfit) to allow us switching storage backends.

The following backends are supported:

* Amazon S3
* Google Cloud Storage
* File system

Usage
=====

It offers you a small API to manipulate your files on multiple storages.

```go
// Storage is the storage interface.
type Storage interface {
	Save(ctx context.Context, content io.Reader, path string) error
	Stat(ctx context.Context, path string) (*Stat, error)
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}

// Stat contains metadata about content stored in storage.
type Stat struct {
	ModifiedTime time.Time
	Size         int64
}

// ErrNotExist is a sentinel error returned by the Stat method.
var ErrNotExist = errors.New("does not exist")
```

If you are migrating from a File system storage to an Amazon S3, you don't need to migrate all your methods anymore!

Be lazy again!

File system
-----------

The file system backend requires a root directory.

```go
dir := os.TempDir()
storage := fs.NewStorage(fs.Config{Root: dir})

// Saving a file named test
storage.Save(ctx, strings.NewReader("(╯°□°）╯︵ ┻━┻"), "test")

// Deleting the new file on the storage
storage.Delete(ctx, "test")
```


Amazon S3
---------

The S3 backend requires AWS credentials, an AWS region and a S3 bucket name.

```go
storage, _ := s3.NewStorage(s3.Config{
    AccessKeyID:     accessKeyID,
    SecretAccessKey: secretAccessKey,
    Region:          region,
    Bucket:          bucket,
})

// Saving a file named test
storage.Save(ctx, strings.NewReader("(>_<)"), "test")

// Deleting the new file on the storage
storage.Delete(ctx, "test")
```

Roadmap
=======

See [issues](https://github.com/ulule/gostorages/issues).

Don't hesitate to send patch or improvements.