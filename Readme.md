# ProGen
![test](https://github.com/kozmod/progen/actions/workflows/test.yml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kozmod/progen)

Simple project's generator.

### Installation

```console
go install github.com/kozmod/progen@latest
```

### About

`progen` use `yml` config file to generate directories, files and execute commands

#### Allowed config file's keys

| Name |  Type  |     Description     |
|:-----|:------:|:-------------------:|
| f    | string | path to config file |
| v    |  bool  |   verbose output    |
| help |  bool  |        help         |

#### Allowed config file's keys

| Key               |       Type        |                      Description                      |
|:------------------|:-----------------:|:-----------------------------------------------------:|
| dirs              |   string slice    |             list of directories to create             |
| files             |      struct       |             list file's `path` and `data`             |
| files.path        |      string       |                       list file                       |
| files.data        |      string       |                       file data                       |
| files.get         |      struct       | struct describe `GET` request for getting file's data |
| files.get.url     |      string       |                      request URL                      |
| files.get.headers | map[string]string |                    request headers                    |
| cmd               |   string slice    |            list of directories to execute             |

<b>Note</b>: preprocessing of "raw" config use [text/template](https://pkg.go.dev/text/template) package
that allow to add custom `yaml` keys tree to avoid duplication (all tags could be used as template's value)

#### Example

```yaml
# custom variables to avoid duplication ( for example "{{.vars.GOPROXY}}")
vars:
  GOPROXY: https://127.0.0.1:8081
  TOKEN: PRIVATE-TOKEN:token
  REPO_1: https://gitlab.some.com/api/v4/projects/23/repository/files

# list directories to create
dirs:
  - api
  - internal/client
  - pkg

# list files to create
files:
  - path: .gitlab-ci.yml
    # GET file from remote storage
    get:
      url: "https://some_file_server.com/files/.gitlab-ci.yml"
      headers:
        some_header: header
  - path: .gitignore
    data: |
      .DS_Store
      .vs/
      .vscode/
      .idea/
      tmp/
  - path: deploy/Dockerfile
    data: |
      FROM golang:1.18.3-alpine as builder

      ENV GOPROXY "{{.vars.GOPROXY}} ,proxy.golang.org,direct"
      ENV GO111MODULE on
      ENV CGO_ENABLED 1
      ENV GOOS linux
      ENV GOARCH amd64

      WORKDIR /app
      COPY . .
      RUN --mount=type=cache,target=/go build -o main .

      FROM alpine:3.16
      ARG config_file

      RUN apk --no-cache --update --upgrade add curl

      WORKDIR /app
      COPY configs/${config_file:-config.yaml} configs/config.yaml
      COPY --from=0 /app/main .
      CMD ["./main"]

# list commands to execute
cmd:
  - curl -H {{.vars.TOKEN}} {{.vars.REPO_1}}/.gitignore/raw?ref=master -o .gitignore
```

```
progen -v -f conf.yml
```

generated project structure

```
.
├── api
├── conf.yml
├── deploy
│   └── Dockerfile
├── internal
│   └── client
├── pkg
└── .gitlab-ci.yml

```

