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
make build
```

___

### About

`progen` use `yml` config file to generate directories, files and execute commands ([actions](#Actions))
___

### Args

| Name    |  Type  |                            Description                             |
|:--------|:------:|:------------------------------------------------------------------:|
| f       | string |                        path to config file                         |
| v       |  bool  |                           verbose output                           |
| dr      |  bool  | `dry run` mode <br/>(to verbose output should be combine with`-v`) |
| version |  bool  |                           print version                            |
| help    |  bool  |                             show flags                             |

___

### Actions

| Key                               |       Type        |    Optional    |                           Description                           |
|:----------------------------------|:-----------------:|:--------------:|:---------------------------------------------------------------:|
|                                   |                   |                |                                                                 |
| http                              |                   |       ✅        |                    http client configuration                    |
| http.debug                        |       bool        |       ✅        |                    http client `DEBUG` mode                     |
| http.base_url                     |      string       |       ✅        |                     http client base `URL`                      |
| http.headers                      | map[string]string |       ✅        |               http client base request `Headers`                |
|                                   |                   |                |                                                                 |
| dirs`<unique_suffix>`<sup>1</sup> |     []string      |       ✅        |                  list of directories to create                  |
|                                   |                   |                |                                                                 |
| files`<unique_suffix>`            |                   |       ✅        |                  list file's `path` and `data`                  |
| files.path                        |      string       |       ❌        |                        save file `path`                         |
| files.tmpl_skip                   |       bool        |       ✅        | flag to skip processing file data as template(except of `data`) |
| files.local                       |      string       | ✳️<sup>2</sup> |                     local file path to copy                     |
| files.data                        |      string       |       ✳️       |                        save file `data`                         |
|                                   |                   |                |                                                                 |
| files.get                         |                   |       ✳️       |      struct describe `GET` request for getting file's data      |
| files.get.url                     |      string       |       ❌        |                          request `URL`                          |
| files.get.headers                 | map[string]string |       ✅        |                         request headers                         |
|                                   |                   |                |                                                                 |
| cmd`<unique_suffix>`              |      []slice      |       ✅        |                   list of command to execute                    |

1. all action execute on declaration order. Base actions (`dir`, `files`,`cmd`) could be configured
   with `<unique_suffix>` to separate action execution.
2. ✳️ only one must be specified in parent section

___

### Usage

#### Generate

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

#### Execution order
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

#### Templates

Configuration preprocessing uses [text/template](https://pkg.go.dev/text/template) of golang's stdlib.
Using templates could be useful to avoiding duplication in configuration file. 
All `text/template` variables must be declared as comments and can be used only to configure data of configuration file (all ones skipping for `file.data` section).
Configuration's `yaml` tag tree also use as `text/template` variables dictionary and can be use for avoiding duplication in configuration file 
and files contents (`files` section).

```yaml
## progen.yml

## `text/template` variables declaration 👇
# {{$project_name := "SOME_PROJECT"}}

## unmapped section (not `dirs`, `files`, `cmd`, `http`) can be use as template variables
vars:
  file_path: some/file/path

dirs:
  - api/{{$project_name}}/v1 # using `text/template` variables
  - internal/{{.vars.file_path}} # using `vars` section
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
#### Http Client
HTTP client configuration
```yaml
## progen.yml

http:
  debug: false
  base_url: https://gitlab.repo_2.com/api/v4/projects/5/repository/files/
  headers:
    PRIVATE-TOKEN: glpat-SOME_TOKEN
```
#### Files
File's content can be declared in configuration file (`files.data` tag) or 
can be received from local file  (`files.local`) or remote (`files.get`). 
Any file's content uses as [text/template](https://pkg.go.dev/text/template) 
and configuration's `yaml` tag tree applies as template variables.

```yaml
## progen.yml

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