package main

import (
  "context"
  "github.com/labstack/echo"
  "github.com/rs/xid"
  "google.golang.org/appengine"
//"google.golang.org/appengine/datastore"
//"google.golang.org/appengine/memcache"
  "cloud.google.com/go/storage"
  "log"
  "net/http"
  "time"
)

const bucketName = "teraconn_raw_voice"

func init() {
  e := echo.New()
  e.GET("/raw_voice_sign", rawVoiceSign)
  e.GET("/lesson_materials/:id", getLessonMaterials)
  e.POST("/lesson_materials", createLessonMaterials)
  http.Handle("/", e)
}

func rawVoiceSign(c echo.Context) error {
  ctx      := c.Request().Context()
  fileID   := xid.New().String()
  lessonID := c.QueryParam("lesson_id")
  fileName := lessonID + "-" + fileID + ".wav"

  createBlankObjectToGCS(ctx, fileName)

  return c.JSON(http.StatusOK, signedURL(ctx, fileID, fileName))
}

func getLessonMaterials(c echo.Context) error {
  // increment view cont in memorycache 
  // https://cloud.google.com/appengine/docs/standard/go/memcache/reference
  // https://cloud.google.com/appengine/docs/standard/go/memcache/using?hl=ja
  id := c.Param("id")
  return c.String(http.StatusOK, id)
}

func createLessonMaterials(c echo.Context) error {
  return c.String(http.StatusCreated, "")
}

func createBlankObjectToGCS(ctx context.Context, fileName string) {
  client, clientErr := storage.NewClient(ctx)
  if clientErr != nil {
    log.Fatal(clientErr)
  }

  bucket := client.Bucket(bucketName)
  obj    := bucket.Object(fileName)
  w      := obj.NewWriter(ctx)

  w.ContentType = "audio/wav"

  if writerErr := w.Close(); writerErr != nil {
    log.Fatal(writerErr)
  }
}

func signedURL(ctx context.Context, fileID string, fileName string) RawWavSign {
  account, _ := appengine.ServiceAccount(ctx)
  expire     := time.Now().AddDate(1, 0, 0)

  url, sign_err := storage.SignedURL(bucketName, fileName, &storage.SignedURLOptions {
    GoogleAccessID: account,
    SignBytes: func(b []byte) ([]byte, error) {
      _, signedBytes, err := appengine.SignBytes(ctx, b)
      return signedBytes, err
	  },
    Method: "PUT",
    Expires: expire,
  })

  if sign_err != nil {
    log.Fatal(sign_err)
  }

  return RawWavSign{ FileID: fileID, SignedURL: url }
}

type RawWavSign struct {
  FileID    string `json:"file_id"`
  SignedURL string `json:"signed_url"`
}
