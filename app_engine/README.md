### setup development environment

```bash
$ brew install goenv
$ goenv install 1.14.10
$ goenv local 1.14.10
$ go get -u
```

### Deploy to App Engine

```bash
$ gcloud components update
$ go get -u github.com/super-dog-human/teraconnectgo

$ gcloud config set project teraconnect-staging
$ gcloud app deploy staging.yaml -v dev0

$ gcloud config set project teraconnect
$ gcloud app deploy production.yaml
```
