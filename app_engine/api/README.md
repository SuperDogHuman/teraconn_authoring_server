### setup development environment
```bash
$ brew install goenv
$ goenv install 1.14.10
$ goenv local 1.14.10
$ go get -u
```

### deploy in local
```bash
$ dev_appserver.py development.yaml --log_level=debug
```

### deploy to Server
```bash
$ gcloud components update

$ gcloud config set project teraconnect-staging
$ gcloud app deploy staging.yaml -v dev0

$ gcloud config set project teraconnect
$ gcloud app deploy production.yaml
```
