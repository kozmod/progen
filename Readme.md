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

#### Action config file's tags

| Key               |       Type        | Optional |                         Description                         |
|:------------------|:-----------------:|:--------:|:-----------------------------------------------------------:|
|                   |                   |          |                                                             |
| http              |                   |    ✅     |                  http client configuration                  |
| http.debug        |       bool        |    ✅     |                  http client `DEBUG` mode                   |
| http.base_url     |      string       |    ✅     |                   http client base `URL`                    |
| http.headers      | map[string]string |    ✅     |             http client base request `Headers`              |
|                   |                   |          |                                                             |
| dirs              |     []string      |    ✅     |                list of directories to create                |
|                   |                   |          |                                                             |
| files             |                   |    ✅     |                list file's `path` and `data`                |
| files.path        |      string       |    ❌     |                      save file `path`                       |
| files.template    |       bool        |    ✅     | flag to apply template variable for file (except of `data`) |
| files.local       |      string       |    ✳️    |                   local file path to copy                   |
| files.data        |      string       |    ✳️    |                      save file `data`                       |
|                   |                   |          |                                                             |
| files.get         |                   |    ✳️    |    struct describe `GET` request for getting file's data    |
| files.get.url     |      string       |    ❌     |                        request `URL`                        |
| files.get.headers | map[string]string |    ✅     |                       request headers                       |
|                   |                   |          |                                                             |
| cmd               |      []slice      |    ✅     |                 list of command to execute                  |

✳️ required one of for parent block

#### Example

```yaml
## preprocessing of "raw" config use `text/template` of golang's stdlib
## and can be used to avoid duplication
## ❗️ `template variables` should be declared 
##     as comments for success `yaml` parsing
##     and only can be applied for current config

## `template variables` declaration 👇🏻 
# {{$gitlab_suffix := "/raw?ref=some_branch"}}

## Example:
## {{printf `%s%s` `.editorconfig` $gitlab_suffix}} ->  .editorconfig/raw?ref=some_branch
## {{$gitlab_suffix}} -> /raw?ref=some_branch

## yaml tags witch not using as `action tags`  👇🏻
## also can be use as `template variables` for current configuration 
## and can be applied to file data got from source (`local`, `get` tags)
vars:
  GOPROXY: https://127.0.0.1:8081
  TOKEN: token
  REPO_1: https://gitlab.repo_1.com/api/v4/projects/23/repository/files
## Example:
## {{.vars.GOPROXY}} -> https://127.0.0.1:8081
## {{.vars.REPO_1}} -> https://gitlab.repo_1.com/api/v4/projects/23/repository/files

# common http client configuration  
http:
  debug: false
  base_url: https://gitlab.repo_2.com/api/v4/projects/5/repository/files/
  headers:
    PRIVATE-TOKEN: {{ .vars.TOKEN }}

# list directories to create 👇🏻
dirs:
  - api
  - internal/client
  - pkg

# list files to create 👇🏻
files:
  - path: .editorconfig
    get:
      url: "{{printf `%s%s` `.editorconfig` $gitlab_suffix}}"

  - path: .gitlab-ci.yml
    # GET file from remote storage
    get:
      # reset url of common http client configuration (http.base_url)
      url: "https://some_file_server.com/files/.gitlab-ci.yml"
      # reset headers of common http client configuration (http.headers)
      headers:
        some_header: header

  - path: Dockerfile
    # process file as template (apply tags which declared in this config)
    template: true
    # GET file from remote storage (using common http client config)
    get:
      # reuse `base` URL of common http client config (http.base_url)
      url: Dockerfile/raw?ref=feature/project_templates"

  - path: .gitignore
    # copy file from location
    local: some/dir/.gitignore.gotmpl

  - path: .env
    # template (false/true) is not necessary - all files with `data` section process as template
    template: true
    data: |
      GOPROXY="{{.vars.GOPROXY}} ,proxy.golang.org,direct"

# list commands to execute 👇🏻
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
├── .env
├── .gitignore
├── .gitlab-ci.yml
└── Dockerfile
```

