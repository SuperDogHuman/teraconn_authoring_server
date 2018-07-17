package main

import (
  "lessonType"
  "context"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/rs/xid"
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
  "google.golang.org/appengine/log"
//"google.golang.org/appengine/memcache"
  "cloud.google.com/go/storage"
  "net/http"
  "time"
)

const bucketName = "teraconn_raw_voice"

func init() {
  e := echo.New()
  e.Pre(middleware.RemoveTrailingSlash())

  e.GET("/raw_voice_signing", rawVoiceSigning)

  e.GET ("/lessons", getLessons)
  e.GET ("/lessons/:id", getLesson)
  e.POST("/lessons", createLesson)
  e.PUT ("/lessons/:id", updateLesson)

  e.GET ("/lessons/:id/materials", getLessonMaterials)
  e.POST("/lessons/:id/materials", createLessonMaterials)
  e.PUT ("/lessons/:id/materials", createLessonMaterials) // same function as POST

  http.Handle("/", e)
}

func rawVoiceSigning(c echo.Context) error {
  ctx      := appengine.NewContext(c.Request())
  fileID   := xid.New().String()
  lessonID := c.QueryParam("lesson_id")
  fileName := lessonID + "-" + fileID + ".wav"

  if err := createBlankObjectToGCS(ctx, fileName); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  if signedURL, err := signedURL(ctx, fileID, fileName); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  } else {
    return c.JSON(http.StatusOK, signedURL)
  }
}

func getLessons(c echo.Context) error {
  // pagination
//  lessons := []lessonType.lesson
  return c.JSON(http.StatusOK, "")
}

func getLesson(c echo.Context) error {
  id       := c.Param("id")
  lesson   := new(lessonType.Lesson)
  lesson.ID = id

  ctx := appengine.NewContext(c.Request())
  if err := fetchLessonFromGCS(ctx, lesson); err != nil {
    if err == datastore.ErrNoSuchEntity {
      log.Errorf(ctx, err.Error())
      return c.JSON(http.StatusNotFound, err.Error())
    } else {
      log.Errorf(ctx, err.Error())
      return c.JSON(http.StatusInternalServerError, err.Error())
    }
  }

  return c.JSON(http.StatusOK, lesson)
}

func createLesson(c echo.Context) error {
  id              := xid.New().String()
  lesson          := new(lessonType.Lesson)
  lesson.ID        = id
  lesson.Published = time.Now()

  ctx := appengine.NewContext(c.Request())
  if err := putLessonToGCS(c, ctx, lesson); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  return c.JSON(http.StatusCreated, lesson)
}

func updateLesson(c echo.Context) error {
  id       := c.Param("id")
  lesson   := new(lessonType.Lesson)
  lesson.ID = id

  ctx := appengine.NewContext(c.Request())
  if err := fetchLessonFromGCS(ctx, lesson); err != nil {
    if err == datastore.ErrNoSuchEntity {
      log.Errorf(ctx, err.Error())
      return c.JSON(http.StatusNotFound, err.Error())
    } else {
      log.Errorf(ctx, err.Error())
      return c.JSON(http.StatusInternalServerError, err.Error())
    }
  }

  if err := putLessonToGCS(c, ctx, lesson); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  return c.JSON(http.StatusOK, lesson)
}

func fetchLessonFromGCS(ctx context.Context, lesson *lessonType.Lesson) error {
  key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)

  if err := datastore.Get(ctx, key, lesson); err != nil {
    return err
  }

  return nil
}

func putLessonToGCS(echoCtx echo.Context, ctx context.Context, lesson *lessonType.Lesson) error {
  if err := echoCtx.Bind(lesson); err != nil {
    return err
  }

  key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)
  lesson.Updated = time.Now()

  if _, err := datastore.Put(ctx, key, lesson); err != nil {
    return err
  }

  return nil
}



// should separate file



func getLessonMaterials(c echo.Context) error {
  // using storage 
  // increment view cont in memorycache 
  // https://cloud.google.com/appengine/docs/standard/go/memcache/reference
  // https://cloud.google.com/appengine/docs/standard/go/memcache/using?hl=ja
  id := c.Param("id")
  return c.JSON(http.StatusOK, id)
}

func createLessonMaterials(c echo.Context) error {
  // using storage 
//  ctx := appengine.NewContext(c.Request())
//  body := c.String(c.Request().Body)
  id := "1"
  return c.JSON(http.StatusCreated, id)
}

func createBlankObjectToGCS(ctx context.Context, fileName string) error {
  client, clientErr := storage.NewClient(ctx)
  if clientErr != nil {
    return clientErr
  }

  bucket := client.Bucket(bucketName)
  obj    := bucket.Object(fileName)
  w      := obj.NewWriter(ctx)

  w.ContentType = "audio/wav"

  if writerErr := w.Close(); writerErr != nil {
    return clientErr
  }

  return nil
}

func signedURL(ctx context.Context, fileID string, fileName string) (RawWavSign, error) {
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
    return RawWavSign{}, sign_err
  }

  return RawWavSign{ FileID: fileID, SignedURL: url }, nil
}

type RawWavSign struct {
  FileID    string `json:"file_id"`
  SignedURL string `json:"signed_url"`
}
