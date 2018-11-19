## production
- gsutil cors set material_cors.json gs:///teraconn_material
- gsutil cors set raw_voice_cors.json gs://teraconn_raw_voice/

## development
- gsutil cors set material_cors_dev.json gs://teraconn_material_development/
- gsutil cors set raw_voice_cors_dev.json gs://teraconn_raw_voice_development/
