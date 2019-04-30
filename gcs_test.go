package gostorages

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGCS(t *testing.T) {
	credentailsFile := ""
	if credentailsFile == "" {
		return
	}
	baseURL := "http://aa"
	storage, err := NewGCSStorage(credentailsFile, "someBucket", "testfolder", baseURL)
	if err != nil {
		t.Fatal(err)
	}

	/*
		filename := "test.txt"
		assert.True(t, storage.Exists(filename))

		storageFile, err := storage.Open(filename)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(5), storageFile.Size())
	*/

	err = storage.Save("test", NewContentFile([]byte("a content example")))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, storage.URL("test"), baseURL+"/testfolder/test")

	modified, err := storage.ModifiedTime("test")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modified.Format(time.RFC822), time.Now().Format(time.RFC822))

	err = storage.Delete("test")

	if err != nil {
		t.Fatal(err)
	}
}
