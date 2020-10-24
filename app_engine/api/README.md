### setup development environment
```bash
$ brew install goenv
$ goenv install 1.14.10
$ goenv local 1.14.10
$ go get -u
```

### deploy in local
```bash
$ dev_appserver.py development/app.yaml --log_level=debug
```

### deploy to Server
```bash
$ gcloud components update
$ gcloud app deploy staging/app.yaml -v dev0
$ gcloud app deploy production/app.yaml
```
