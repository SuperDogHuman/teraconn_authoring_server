package main

import (
  "lessonType"
// "lesson"
// "lessonMaterial"
// "rawVoiceSigning"
//  "GCSHelper"
  "bytes"
  "context"
  "encoding/json"
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

func init() {
  e := echo.New()
  e.Pre(middleware.RemoveTrailingSlash())
  e.Use(middleware.CORS())

  e.GET("/raw_voice_signing", rawVoiceSigning)

  e.GET ("/lessons", getLessons)
  e.GET ("/lessons/:id", getLesson)
  e.POST("/lessons", createLesson)
  e.PUT ("/lessons/:id", updateLesson)

  e.GET ("/lessons/:id/materials", getLessonMaterials)
  e.POST("/lessons/:id/materials", putLessonMaterials)
  e.PUT ("/lessons/:id/materials", putLessonMaterials) // same function as POST

  http.Handle("/", e)
}

func rawVoiceSigning(c echo.Context) error {
  ctx         := appengine.NewContext(c.Request())
  bucketName  := "teraconn_raw_voice"
  fileID      := xid.New().String()
  lessonID    := c.QueryParam("lesson_id")
  fileName    := lessonID + "-" + fileID + ".wav"
  contentType := "audio/wav"

  if err := createObjectToGCS(ctx, bucketName, fileName, contentType, nil); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  if signedURL, err := signedURL(ctx, bucketName, fileID, fileName); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  } else {
    return c.JSON(http.StatusOK, signedURL)
  }
}

func voiceText(c echo.Context) error {
  // [ids]
  // return id, voicetext
  return c.JSON(http.StatusOK, "")
}

func getLessons(c echo.Context) error {
  // pagination
  return c.JSON(http.StatusOK, "")
}

func getLesson(c echo.Context) error {
  id       := c.Param("id")
  lesson   := new(lessonType.Lesson)
  lesson.ID = id

  ctx := appengine.NewContext(c.Request())
  if err := fetchLessonFromGCD(ctx, lesson); err != nil {
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
  if err := putLessonToGCD(c, ctx, lesson); err != nil {
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
  if err := fetchLessonFromGCD(ctx, lesson); err != nil {
    if err == datastore.ErrNoSuchEntity {
      log.Errorf(ctx, err.Error())
      return c.JSON(http.StatusNotFound, err.Error())
    } else {
      log.Errorf(ctx, err.Error())
      return c.JSON(http.StatusInternalServerError, err.Error())
    }
  }

  if err := putLessonToGCD(c, ctx, lesson); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  return c.JSON(http.StatusOK, lesson)
}

func fetchLessonFromGCD(ctx context.Context, lesson *lessonType.Lesson) error {
  key := datastore.NewKey(ctx, "Lesson", lesson.ID, 0, nil)

  if err := datastore.Get(ctx, key, lesson); err != nil {
    return err
  }

  return nil
}

func putLessonToGCD(echoCtx echo.Context, ctx context.Context, lesson *lessonType.Lesson) error {
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

func getLessonMaterials(c echo.Context) error {
  // increment view cont in memorycache 
  // https://cloud.google.com/appengine/docs/standard/go/memcache/reference
  // https://cloud.google.com/appengine/docs/standard/go/memcache/using?hl=ja

  lessonID    := c.Param("id")
  ctx         := appengine.NewContext(c.Request())
  bucketName  := "teraconn_material"
  filePath    := "lesson/" + lessonID + ".json"

  bytes, err := getObjectFromGCS(ctx, bucketName, filePath)
  if err != nil {
    log.Errorf(ctx, err.Error())
    if err == storage.ErrObjectNotExist {
      return c.JSON(http.StatusNotFound, err.Error())
    } else {
      return c.JSON(http.StatusInternalServerError, err.Error())
    }
  }

  lessonMaterial := new(lessonType.LessonMaterial)
  if err := json.Unmarshal(bytes, lessonMaterial); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  return c.JSON(http.StatusOK, lessonMaterial)
}

func putLessonMaterials(c echo.Context) error {
  lessonID       := c.Param("id")
  ctx            := appengine.NewContext(c.Request())
  lessonMaterial := new(lessonType.LessonMaterial)

  if err := c.Bind(lessonMaterial); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  contents, err := json.Marshal(lessonMaterial)
  if err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  bucketName  := "teraconn_material"
  filePath    := "lesson/" + lessonID + ".json"
  contentType := "application/json"

  if err := createObjectToGCS(ctx, bucketName, filePath, contentType, contents); err != nil {
    log.Errorf(ctx, err.Error())
    return c.JSON(http.StatusInternalServerError, err.Error())
  }

  return c.JSON(http.StatusCreated, "succeed")
}

func createObjectToGCS(ctx context.Context, bucketName, filePath, contentType string, contents []byte) error {
  client, err := storage.NewClient(ctx)
  if err != nil { return err }
  defer client.Close()

  w := client.Bucket(bucketName).Object(filePath).NewWriter(ctx)
  w.ContentType = contentType
  defer w.Close()

  if (len(contents) > 0) {
    if _, err := w.Write(contents); err != nil {
      return err
    }
  }

  if err := w.Close(); err != nil {
    return err
  }

  return nil
}

func getObjectFromGCS(ctx context.Context, bucketName, filePath string) ([]byte, error) {
  client, err := storage.NewClient(ctx)
  if err != nil { return nil, err }
  defer client.Close()

  r, err := client.Bucket(bucketName).Object(filePath).NewReader(ctx)
  if err != nil { return nil, err }
  defer r.Close()

  var buffer bytes.Buffer
  if _, err := buffer.ReadFrom(r); err != nil {
      return nil, err
  }

  return buffer.Bytes(), nil
}

func signedURL(ctx context.Context, bucketName string, fileID string, fileName string) (RawWavSign, error) {
  account, _ := appengine.ServiceAccount(ctx)
  expire     := time.Now().AddDate(1, 0, 0)

  url, sign_err := storage.SignedURL(bucketName, fileName, &storage.SignedURLOptions {
    GoogleAccessID: account,
    SignBytes: func(b []byte) ([]byte, error) {
      _, signedBytes, err := appengine.SignBytes(ctx, b)
      return signedBytes, err
	  },
    Method: "PUT",
    ContentType: "audio/wav",
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
