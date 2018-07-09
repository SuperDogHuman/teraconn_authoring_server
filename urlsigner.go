package urlsigner 

import (
  "bytes"
  "context"
  "encoding/json"
  "google.golang.org/appengine"
  "google.golang.org/appengine/file"
  "google.golang.org/cloud/storage"
  "log"
  "net/http"
  "time"
)

var bucket string
var ctx context.Context
var filename string

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx = appengine.NewContext(r)
    bucket, _ = file.DefaultBucketName(ctx)
    filename = r.URL.Query()["filename"][0]

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(signedURL(w))
}

func signedURL(w http.ResponseWriter) []byte {
  account, _ := appengine.ServiceAccount(ctx)
  expire := time.Now().AddDate(1, 0, 0)

  url, sign_err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions {
    GoogleAccessID: account,
    SignBytes: func(b []byte) ([]byte, error) {
      _, signedBytes, err := appengine.SignBytes(ctx, b)
      return signedBytes, err
	  },
    Method: "GET",
    Expires: expire,
  })

  if sign_err != nil {
    log.Fatal(sign_err)
  }

  response := Response{ SignedURL: url }

  buffer := bytes.NewBuffer([]byte{})
  jsonEncoder := json.NewEncoder(buffer)
  jsonEncoder.SetEscapeHTML(false)
  jsonEncoder.Encode(response)

  return buffer.Bytes()
}

type Response struct {
  SignedURL string `json:"signed_url"`
}
