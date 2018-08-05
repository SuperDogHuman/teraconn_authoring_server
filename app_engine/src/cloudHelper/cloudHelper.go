package cloudHelper

import (
	"bytes"
	"context"
	"lessonType"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// FetchObjectFromGCD is fetch object from GCD function.
func FetchObjectFromGCD(ctx context.Context, obj *lessonType.Lesson) error {
	key := datastore.NewKey(ctx, "Lesson", obj.ID, 0, nil)

	if err := datastore.Get(ctx, key, obj); err != nil {
		return err
	}

	return nil
}

// PutObjectToGCD is fetch object from GCD function.
func PutObjectToGCD(ctx context.Context, echoCtx echo.Context, obj *lessonType.Lesson) error {
	if err := echoCtx.Bind(obj); err != nil {
		return err
	}

	key := datastore.NewKey(ctx, "Lesson", obj.ID, 0, nil)
	obj.Updated = time.Now()

	if _, err := datastore.Put(ctx, key, obj); err != nil {
		return err
	}

	return nil
}

// CreateObjectToGCS is create object to GCS function.
func CreateObjectToGCS(ctx context.Context, bucketName, filePath, contentType string, contents []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	w := client.Bucket(bucketName).Object(filePath).NewWriter(ctx)
	w.ContentType = contentType
	defer w.Close()

	if len(contents) > 0 {
		if _, err := w.Write(contents); err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

// GetObjectFromGCS is get object from GCS function.
func GetObjectFromGCS(ctx context.Context, bucketName, filePath string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	r, err := client.Bucket(bucketName).Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(r); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// GetGCSSignedURL is generate signed-URL for GCS object.
func GetGCSSignedURL(ctx context.Context, bucketName string, fileID string, fileName string, method string, contentType string) (string, error) {
	account, _ := appengine.ServiceAccount(ctx)
	expire := time.Now().AddDate(1, 0, 0)

	url, signErr := storage.SignedURL(bucketName, fileName, &storage.SignedURLOptions{
		GoogleAccessID: account,
		SignBytes: func(b []byte) ([]byte, error) {
			_, signedBytes, err := appengine.SignBytes(ctx, b)
			return signedBytes, err
		},
		Method:      method,
		ContentType: contentType,
		Expires:     expire,
	})

	if signErr != nil {
		return url, signErr
	}

	return url, nil
}
