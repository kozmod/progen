# ProGen

![test](https://github.com/kozmod/progen/actions/workflows/test.yml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kozmod/progen)

Simple projects generator.

### Installation

```console
go install github.com/kozmod/progen@latest
```

### Build from source
```console
go build -o progen .
```

### About

`progen` use `yml` config file to generate directories, files and execute commands

#### Allowed config file's keys

| Name |  Type  |     Description     |
|:-----|:------:|:-------------------:|
| f    | string | path to config file |
| v    |  bool  |   verbose output    |
| help |  bool  |  flags information  |

#### Allowed config file's keys

| Key               |       Type        | Optional |                         Description                         |
|:------------------|:-----------------:|:--------:|:-----------------------------------------------------------:|
|                   |                   |          |                                                             |
| http              |      struct       |    ✅     |                  http client configuration                  |
| http.debug        |       bool        |    ✅     |                  http client `DEBUG` mode                   |
| http.base_url     |      string       |    ✅     |                   http client base `URL`                    |
| http.headers      | map[string]string |    ✅     |             http client base request `Headers`              |
|                   |                   |          |                                                             |
| dirs              |     []string      |    ✅     |                list of directories to create                |
|                   |                   |          |                                                             |
| files             |      struct       |    ✅     |                list file's `path` and `data`                |
| files.path        |      string       |    ❌     |                      save file `path`                       |
| files.template    |       bool        |    ✅     | flag to apply template variable for file (except of `data`) |
| files.local       |      string       |    ✳️    |                   local file path to copy                   |
| files.data        |      string       |    ✳️    |                      save file `data`                       |
| files.get         |      struct       |    ✳️    |    struct describe `GET` request for getting file's data    |
| files.get.url     |      string       |    ❌     |                        request `URL`                        |
| files.get.headers | map[string]string |    ✅     |                       request headers                       |
|                   |                   |          |                                                             |
| cmd               |      []slice      |    ✅     |                 list of command to execute                  |

✳️ required one of for parent block

❗️<b>Note</b>: preprocessing of "raw" config use [text/template](https://pkg.go.dev/text/template) package
that allow to add custom `yaml` keys tree to avoid duplication (all tags could be used as template's value)

#### Example

```yaml
# custom variables to avoid duplication ( for example "{{.vars.GOPROXY}}")
vars:
  GOPROXY: https://127.0.0.1:8081
  TOKEN: token
  REPO_1: https://gitlab.repo_1.com/api/v4/projects/23/repository/files

# common http client configuration  
http:
  debug: false
  base_url: https://gitlab.repo_2.com/api/v4/projects/5/repository/files/
  headers:
    PRIVATE-TOKEN: {{ .vars.TOKEN }}

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
      # reset url of common http client configuration 
      url: "https://some_file_server.com/files/.gitlab-ci.yml"
      # reset headers of common http client configuration (tag:http)
      headers:
        some_header: header

  - path: Dockerfile
    # process file as template (apply variables which declared in this config)
    template: true
    # GET file from remote storage (using common http client config)
    get:
      # reuse `base` URL of common http client config
      url: Dockerfile/raw?ref=feature/project_templates"

  - path: .gitignore
    # copy file from location
    local: some/dir/.gitignore.gotmpl

  - path: .env
    # template (false/true) is not necessary - all files with `data` section process as template
    template: true
    data: |
      GOPROXY="{{.vars.GOPROXY}} ,proxy.golang.org,direct"

# list commands to execute
cmd:
  - curl -H PRIVATE-TOKEN:{{.vars.TOKEN}} {{.vars.REPO_1}}/.editorconfig/raw?ref=master -o .editorconfig
```

#### Generate project structure from configuration file

use configuration file with default name (`progen.yaml`)

```console
progen -v
```

or define custom config location using `-f`

```console
progen -v -f conf.yml
```

generated project structure

```console
.
├── api
├── progen.yml
├── internal
│   └── client
├── pkg
├── .editorconfig 
├── .gitlab-ci.yml
├── .env
└── Dockerfile
```

