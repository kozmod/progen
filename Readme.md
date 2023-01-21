# ProGen

![test](https://github.com/kozmod/progen/actions/workflows/test.yml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kozmod/progen)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/kozmod/progen)

Simple projects generator.
___

### Installation

```console
go install github.com/kozmod/progen@latest
```

### Build from source

```console
go build -o progen .
```

___

### About

`progen` use `yml` config file to generate directories, files and execute commands ([actions](#Actions))
___

### Args

| Name |  Type  |                            Description                             |
|:-----|:------:|:------------------------------------------------------------------:|
| f    | string |                        path to config file                         |
| v    |  bool  |                           verbose output                           |
| dr   |  bool  | `dry run` mode <br/>(to verbose output should be combine with`-v`) |
| help |  bool  |                             show flags                             |

___

### Actions

| Key                               |       Type        |   Optional    |                           Description                           |
|:----------------------------------|:-----------------:|:-------------:|:---------------------------------------------------------------:|
|                                   |                   |               |                                                                 |
| http                              |                   |       ✅       |                    http client configuration                    |
| http.debug                        |       bool        |       ✅       |                    http client `DEBUG` mode                     |
| http.base_url                     |      string       |       ✅       |                     http client base `URL`                      |
| http.headers                      | map[string]string |       ✅       |               http client base request `Headers`                |
|                                   |                   |               |                                                                 |
| dirs`<unique_suffix>`<sup>1</sup> |     []string      |       ✅       |                  list of directories to create                  |
|                                   |                   |               |                                                                 |
| files`<unique_suffix>`            |                   |       ✅       |                  list file's `path` and `data`                  |
| files.path                        |      string       |       ❌       |                        save file `path`                         |
| files.tmpl_skip                   |       bool        |       ✅       | flag to skip processing file data as template(except of `data`) |
| files.local                       |      string       | ✳<sup>2</sup> |                     local file path to copy                     |
| files.data                        |      string       |       ✳       |                        save file `data`                         |
|                                   |                   |               |                                                                 |
| files.get                         |                   |       ✳       |      struct describe `GET` request for getting file's data      |
| files.get.url                     |      string       |       ❌       |                          request `URL`                          |
| files.get.headers                 | map[string]string |       ✅       |                         request headers                         |
|                                   |                   |               |                                                                 |
| cmd`<unique_suffix>`              |      []slice      |       ✅       |                   list of command to execute                    |

1. all action execute on declaration order. Base actions (`dir`, `files`,`cmd`) could be configured
   with `<unique_suffix>` to separate action execution.
2. `✳` only one must be specified in parent section

___

### Usage

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
    PRIVATE-TOKEN: { { .vars.TOKEN } }

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
    tmpl_skip: false
    # GET file from remote storage (using common http client config)
    get:
      # reuse `base` URL of common http client config (http.base_url)
      url: Dockerfile/raw?ref=feature/project_templates"

  - path: .gitignore
    # copy file from location
    local: some/dir/.gitignore.gotmpl

  - path: .env
    # skip file processing as template
    tmpl_skip: true
    data: |
      GOPROXY="{{.vars.GOPROXY}} ,proxy.golang.org,direct"

# list commands to execute (with zero `<unique_suffix>`)👇🏻    
cmd:
  - pwd

# list commands to execute (with `_2` `<unique_suffix>`) 👇🏻
cmd_2:
  - curl -H PRIVATE-TOKEN:{{.vars.TOKEN}} {{.vars.REPO_1}}/buf.gen.yaml/raw?ref=master -o buf.gen.yaml
```

___

### Generate

`progen` use `progen.yaml` as default configuration file

```console
progen -v
```

`-f` flag set custom configuration file

```console
progen -v -f conf.yml
```

generated files and directories

```console
.
├── api
├── progen.yml
├── internal
│   └── client
├── pkg
├── buf.gen.yaml
├── .editorconfig 
├── .env
├── .gitignore
├── .gitlab-ci.yml
└── Dockerfile
```

