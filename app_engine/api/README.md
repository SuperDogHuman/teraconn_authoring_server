- brew install goenv direnv
- go get -u necessary-packages
- cd ./app_engine/api
- gcloud app deploy src/main/development.yaml -v dev0
- gcloud app deploy src/main/production.yaml

