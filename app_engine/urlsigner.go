package urlsigner

import (
  "bytes"
  "context"
  "encoding/json"
  "github.com/rs/xid"
  "google.golang.org/appengine"
  "google.golang.org/cloud/storage"
  "log"
  "net/http"
  "time"
)

var ctx context.Context
var fileID string
var fileName string
const bucketName = "teraconn_raw_voice"

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx           = appengine.NewContext(r)
    fileID        = xid.New().String()
    lessonID     := r.URL.Query().Get("lesson_id")
    fileName      = lessonID + "-" + fileID + ".wav"

    createBlankObject()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(signedURL())
}

func createBlankObject() {
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

func signedURL() []byte {
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

  response := Response{ SignedURL: url, FileID: fileID }
  return responceToJson(response)
}

func responceToJson(response Response) []byte{
  buffer := bytes.NewBuffer([]byte{})

  jsonEncoder := json.NewEncoder(buffer)
  jsonEncoder.SetEscapeHTML(false)
  jsonEncoder.Encode(response)

  return buffer.Bytes()
}

type Response struct {
  FileID    string `json:"file_id"`
  SignedURL string `json:"signed_url"`
}
