# ProGen <img align="right" src=".github/assets/PG1-4-3-1.png" alt="drawing"  width="60" />

![test](https://github.com/kozmod/progen/actions/workflows/test.yml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kozmod/progen)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/kozmod/progen)
![GitHub release date](https://img.shields.io/github/release-date/kozmod/progen)
![GitHub last commit](https://img.shields.io/github/last-commit/kozmod/progen)
![GitHub MIT license](https://img.shields.io/github/license/kozmod/progen)

Simple projects generator.
___

### Installation

```shell
go install github.com/kozmod/progen@latest
```

### Build from source

```shell
make build
```

___

### About

`progen` use `yml` config file to generate directories, files and execute commands ([actions](#Actions))
___

### Args

| Name       |  Type  | Description                                                                                  |
|:-----------|:------:|:---------------------------------------------------------------------------------------------|
| `-f`       | string | path to config file                                                                          |
| `-v`       |  bool  | verbose output                                                                               |
| `-dr`      |  bool  | `dry run` mode <br/>(to verbose output should be combine with`-v`)                           |
| `-awd`     |  bool  | application working directory                                                                |
| `-tvar`    |  bool  | [text/template](https://pkg.go.dev/text/template) variables (override config variables tree) |
| `-version` |  bool  | print version                                                                                |
| `-help`    |  bool  | show flags                                                                                   |

___

### Actions

| Key                                                   |       Type        | Optional | Description                                                     |
|:------------------------------------------------------|:-----------------:|:---------|:----------------------------------------------------------------|
|                                                       |                   |          |                                                                 |
| settings                                              |                   | ✅        | `progen` settings section                                       |
|                                                       |                   |          |
| settings.http                                         |                   | ✅        | http client configuration                                       |
| settings.http.debug                                   |       bool        | ✅        | http client `DEBUG` mode                                        |
| settings.http.base_url                                |      string       | ✅        | http client base `URL`                                          |
| settings.http.headers                                 | map[string]string | ✅        | http client base request `Headers`                              |
| settings.http.query_params                            | map[string]string | ✅        | http client base request `Query Parameters`                     |
|                                                       |                   |          |                                                                 |
| dirs`<unique_suffix>`[<sup>1</sup>](#Execution order) |     []string      | ✅        | list of directories to create                                   |
|                                                       |                   |          |                                                                 |
| files`<unique_suffix>`                                |                   | ✅        | list file's `path` and `data`                                   |
| files.path                                            |      string       | ❌        | save file `path`                                                |
| files.tmpl_skip                                       |       bool        | ✅        | flag to skip processing file data as template(except of `data`) |
| files.local                                           |      string       | `❕`      | local file path to copy                                         |
| files.data                                            |      string       | `❕`      | save file `data`                                                |
|                                                       |                   |          |                                                                 |
| files.get                                             |                   | `❕`      | struct describe `GET` request for getting file's data           |
| files.get.url                                         |      string       | ❌        | request `URL`                                                   |
| files.get.headers                                     | map[string]string | ✅        | request `Headers`                                               |
| files.query_params                                    | map[string]string | ✅        | request `Query Parameters`                                      |
|                                                       |                   |          |                                                                 |
| cmd`<unique_suffix>`                                  |      []slice      | ✅        | list of command to execute                                      |

`❕` only one must be specified in parent section

___

## Usage

### Generate

`prohen` execute commands and generate files and directories based on configuration file

```yaml
## progen.yml

# list directories to creation
dirs:
  - x/y

# list files to creation
files:
  - path: x/some_file.txt
    data: |
      some data

# list commands to execution 
cmd:
  - touch second_file.txt
  - tree
```

```console
% progen -v
2023-01-22 12:44:55	INFO	dir created: x/y
2023-01-22 12:44:55	INFO	file created [template: false]: x/some_file.txt
2023-01-22 12:44:55	INFO	execute: touch second_file.txt
2023-01-22 12:44:55	INFO	execute: tree
out:
.
├── progen.yml
├── second_file.txt
└── x
    ├── some_file.txt
    └── y
```

### Execution order

All actions execute in declared order. Base actions (`dir`, `files`,`cmd`) could be configured
with `<unique_suffix>` to separate action execution.

```yaml
## progen.yml

dirs1:
  - api/some_project/v1
cmd1:
  - chmod -R 777 api
dirs2:
  - api/some_project_2/v1
cmd2:
  - chmod -R 777 api
```

```
% progen -v
2023-01-22 13:38:52	INFO	dir created: api/some_project/v1
2023-01-22 13:38:52	INFO	execute: chmod -R 777 api
2023-01-22 13:38:52	INFO	dir created: api/some_project_2/v1
2023-01-22 13:38:52	INFO	execute: chmod -R 777 api
```

### Templates

Configuration preprocessing uses [text/template](https://pkg.go.dev/text/template) of golang's stdlib.
Using templates could be useful to avoiding duplication in configuration file.
All `text/template` variables must be declared as comments and can be used only to configure data of configuration
file (all ones skipping for `file.data` section).
Configuration's `yaml` tag tree also use as `text/template` variables dictionary and can be use for avoiding duplication
in configuration file
and files contents (`files` section).

```yaml
## progen.yml

## `text/template` variables declaration 👇
# {{$project_name := "SOME_PROJECT"}}

## unmapped section (not `dirs`, `files`, `cmd`, `http`) can be use as template variables
vars:
  file_path: some/file/path

dirs:
  - api/{{$project_name}}/v1 # used from `text/template` variables
  - internal/{{.vars.file_path}} # used from `vars` section
  - pkg/{{printf `%s-%s` $project_name `data`}}

files:
  - path: internal/{{$project_name}}.txt
    data: |
      Project name:{{$project_name}}
  - path: pkg/{{printf `%s-%s` $project_name `data`}}/some_file.txt
    tmpl_skip: true
    data: |
      {{$project_name}}

cmd:
  - cat internal/{{$project_name}}.txt
  - cat pkg/{{printf `%s-%s` $project_name `data`}}/some_file.txt
  - tree
```

```console
% progen -v
2023-01-22 13:03:58	INFO	dir created: api/SOME_PROJECT/v1
2023-01-22 13:03:58	INFO	dir created: internal/some/file/path
2023-01-22 13:03:58	INFO	dir created: pkg/SOME_PROJECT-data
2023-01-22 13:03:58	INFO	file created [template: true]: internal/SOME_PROJECT.txt
2023-01-22 13:03:58	INFO	file created [template: false]: pkg/SOME_PROJECT-data/some_file.txt
2023-01-22 13:03:58	INFO	execute: cat internal/SOME_PROJECT.txt
out:
Project name:SOME_PROJECT

2023-01-22 13:03:58	INFO	execute: cat pkg/SOME_PROJECT-data/some_file.txt
out:
{{$project_name}}

2023-01-22 13:03:58	INFO	execute: tree
out:
.
├── api
│   └── SOME_PROJECT
│       └── v1
├── internal
│   ├── SOME_PROJECT.txt
│   └── some
│       └── file
│           └── path
├── pkg
│   └── SOME_PROJECT-data
│       └── some_file.txt
└── progen.yml
```

any part of template variable tree can be override using `-tvar` flag

```yaml
## `text/template` variables declaration 👇
# {{$project_name := "SOME_PROJECT"}}

## unmapped section (not `dirs`, `files`, `cmd`, `http`) can be use as template variables
vars:
  file_path: some/file/path
  file_path_2: some/file/path_2

dirs:
  - api/{{$project_name}}/v1 # used from `text/template` variables
  - internal/{{.vars.file_path}} # used from `vars` section
  - internal/{{.vars.file_path_2}} # used overridden `vars` which set through args (-tvar=.vars.file_path 2=override path)
```

```console
% progen -v -dr -tvar=.vars.file_path_2=overrided_path
2023-01-22 22:25:47	INFO	configuration file: progen_test_vars.yml
2023-01-22 22:25:47	INFO	dir created: api/SOME_PROJECT/v1
2023-01-22 22:25:47	INFO	dir created: internal/some/file/path
2023-01-22 22:25:47	INFO	dir created: internal/overrided_path
```

### Text files

Instead of specifying a config file, you can pass a single `progen.yml` in the pipe the file in via `STDIN`.
To pipe a `progen.yml` from `STDIN`:

```console
progen - < progen.yml
```

or

```console
cat progen.yml | progen -
```

If you use `STDIN`  the system ignores any `-f` option.

**Example** (get `progen.yml` from gitlab repository with replacing [text/template](https://pkg.go.dev/text/template)
variables using `-tvar` flag):

```console
curl -H PRIVATE-TOKEN:token https://gitlab.some.com/api/v4/projects/13/repository/files/shared%2Fteplates%2Fsimple%2Fprogen.yml/raw\?ref\=feature/templates | progen -v -dr -tvar=.vars.GOPROXY=some_proxy -
```

### Http Client

HTTP client configuration

```yaml
## progen.yml

settings:
  http:
    debug: false
    base_url: https://gitlab.repo_2.com/api/v4/projects/5/repository/files/
    headers:
      PRIVATE-TOKEN: glpat-SOME_TOKEN
    query_params:
      PARAM_1: Val_1
```

### Files

File's content can be declared in configuration file (`files.data` tag) or
can be received from local file  (`files.local`) or remote (`files.get`).
Any file's content uses as [text/template](https://pkg.go.dev/text/template)
and configuration's `yaml` tag tree applies as template variables.

```yaml
## progen.yml

# settings of the cli
settings:
  # common http client configuration  
  http:
    debug: false
    base_url: https://gitlab.repo_2.com/api/v4/projects/5/repository/files/
    headers:
      PRIVATE-TOKEN: glpat-SOME_TOKEN

# {{$project_name := "SOME_PROJECT"}}
# {{$gitlab_suffix := "/raw?ref=some_branch"}}

files:
  - path: files/Readme.md
    # skip file processing as template
    tmpl_skip: true
    data: |
      Project name: {{$project_name}}

  - path: files/.gitignore
    # copy file from location
    local: some/dir/.gitignore.gotmpl

  - path: files/.editorconfig
    get:
      url: "{{printf `%s%s` `.editorconfig` $gitlab_suffix}}"

  - path: files/.gitlab-ci.yml
    # GET file from remote storage
    get:
      # reset URL which set in http client configuration (http.base_url)
      url: "https://some_file_server.com/files/.gitlab-ci.yml"
      # reset headers of common http client configuration (http.headers)
      headers:
        some_header: header
      query_params:
        PARAM_1: Val_1

  - path: files/Dockerfile
    # process file as template (apply tags which declared in this config)
    tmpl_skip: false
    # GET file from remote storage (using common http client config)
    get:
      # reuse `base` URL of common http client config (http.base_url)
      url: Dockerfile/raw?ref=feature/project_templates"
```

```console
% progen -v
2023-01-22 15:47:45	INFO	file created [template: false]: files/Readme.md
2023-01-22 15:47:45	INFO	file created [template: true]: files/.gitignore
2023-01-22 15:47:45	INFO	file created [template: true]: files/.editorconfig
2023-01-22 15:47:45	INFO	file created [template: true]: files/.gitlab-ci.yml
2023-01-22 15:47:45	INFO	file created [template: true]: files/Dockerfile
```