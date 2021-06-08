## production
- gsutil cors set material_cors.json gs:///teraconn_material
- gsutil cors set public_cors.json gs://teraconn_public

## development
- gsutil cors set material_cors_dev.json gs://teraconn_material_development
- gsutil cors set public_cors_dev.json gs://teraconn_public_development
