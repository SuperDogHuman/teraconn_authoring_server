package urlsigner

import (
  "bytes"
  "context"
  "encoding/json"
  "github.com/rs/xid"
  "google.golang.org/appengine"
  "google.golang.org/appengine/file"
  "google.golang.org/cloud/storage"
  "log"
  "net/http"
  "time"
)

var ctx context.Context
var bucket string

func init() {
    http.HandleFunc("/", handler)
}

// ルーティング
// ファイル名はこちら側で生成
// GETからPUTに変更（何もなしでいきなりPUTできる？
func handler(w http.ResponseWriter, r *http.Request) {
    ctx = appengine.NewContext(r)
    bucket, _ = file.DefaultBucketName(ctx)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(signedURL())
}

func signedURL() []byte {
  account, _ := appengine.ServiceAccount(ctx)
  expire     := time.Now().AddDate(1, 0, 0)
  filename   := xid.New().String() + ".wav"

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

  response := Response{ SignedURL: url, Filename: filename }
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
  Filename string `json:"filename"`
  SignedURL string `json:"signed_url"`
}
