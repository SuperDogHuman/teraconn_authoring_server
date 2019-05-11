- brew install goenv goenv
- goenv install 1.11.9
- goenv local 1.11.9
- go get -u

- deploy in local
  - dev_appserver.py development.yaml --log_level=debug

- deploy to Server
  - gcloud components update
  - gcloud app deploy development.yaml -v dev0
  - gcloud app deploy production.yaml 
