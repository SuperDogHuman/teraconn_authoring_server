### deploy command

```
$ gcloud config set project teraconnect-development
$ gcloud functions deploy Mp3SpeechToText --runtime go113 --trigger-resource teraconn_material_development --trigger-event google.storage.object.finalize --memory 256 --timeout 540 --region asia-northeast1 --set-env-vars PROJECT_ID=teraconnect-development
```
