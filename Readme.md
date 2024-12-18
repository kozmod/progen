# ProGen <img align="right" src=".github/assets/PG1-4-3-1.png" alt="drawing"  width="60" />

[![test](https://github.com/kozmod/progen/actions/workflows/test.yml/badge.svg)](https://github.com/kozmod/progen/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kozmod/progen)](https://goreportcard.com/report/github.com/kozmod/progen)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kozmod/progen)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/kozmod/progen)
![GitHub release date](https://img.shields.io/github/release-date/kozmod/progen)
![GitHub last commit](https://img.shields.io/github/last-commit/kozmod/progen)
[![GitHub MIT license](https://img.shields.io/github/license/kozmod/progen)](https://github.com/kozmod/progen/blob/main/LICENSE)

A flexible, language and frameworks agnostic tool that allows you to generate projects structure from templates based
on `yaml` configuration (generate directories, files and execute commands) or use as library to build custom generator.
___

### Installation

```shell
go install github.com/kozmod/progen@latest
```

### Build from source

```shell
make build
```
### Use as `lib` [<sup>**ⓘ**</sup>](#lib_usage)
```go
module github.com/some/custom_gen

go 1.22

require github.com/kozmod/progen v0.1.8
```

___

### Flags

| Name                                                                  |   Type   |   Default    | Description                                                                                                                                                                            |
|:----------------------------------------------------------------------|:--------:|:------------:|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-f`[<sup>**ⓘ**</sup>](#config_file) <sup>**✱**</sup>                 |  string  | `progen.yml` | specify configuration file path                                                                                                                                                        |
| `-v` <sup>**✱**</sup>                                                 |   bool   |   `false`    | verbose output                                                                                                                                                                         |
| `-dr`[<sup>**ⓘ**</sup>](#dry_run) <sup>**✱**</sup>                    |   bool   |   `false`    | `dry run` mode <br/>(to verbose output should be combine with`-v`)                                                                                                                     |
| `-awd`[<sup>**ⓘ**</sup>](#awd)                                        |  string  |     `.`      | application working directory                                                                                                                                                          |
| `-printconf`[<sup>**ⓘ**</sup>](#print_config)                         |   bool   |   `false`    | output processed config                                                                                                                                                                |
| `-errtrace`[<sup>**ⓘ**</sup>](#print_err_trace) <sup>**✱**</sup>      |   bool   |   `false`    | output errors stack trace                                                                                                                                                              |
| `-pf`[<sup>**ⓘ**</sup>](#files_preprocessing)                         |   bool   |    `true`    | `preprocessing files`: load and process all files <br/>(all files `actions`[<sup>**ⓘ**</sup>](#files_actio_desk)) as [text/template](https://pkg.go.dev/text/template) before creating |
| `-tvar`[<sup>**ⓘ**</sup>](#tvar) <sup>**✱**</sup>                     | []string |    `[ ]`     | [text/template](https://pkg.go.dev/text/template) variables <br/>(override config variables tree)                                                                                      |
| `-missingkey` <sup>**✱**</sup>                                        | []string |   `error`    | set `missingkey`[text/template.Option](https://pkg.go.dev/text/template#Template.Option) execution option                                                                              |
| `-skip`[<sup>**ⓘ**</sup>](#skip_actions)                              | []string |    `[ ]`     | skip any `action` tag <br/>(regular expression)                                                                                                                                        |
| `-gp`[<sup>**ⓘ**</sup>](#groups_of_actions)                           | []string |    `[ ]`     | set of the action's groups to execution                                                                                                                                                |
| `-version`                                                            |   bool   |   `false`    | print version                                                                                                                                                                          |
| `-help` <sup>**✱**</sup>                                              |   bool   |   `false`    | show flags                                                                                                                                                                             |

<sup>**✱**</sup> flags accessible in the `lib`. 
___

### Actions and tags

| Key                                                                             |       Type        | Optional | Description                                                                                                 |
|:--------------------------------------------------------------------------------|:-----------------:|:---------|:------------------------------------------------------------------------------------------------------------|
|                                                                                 |                   |          |                                                                                                             |
| settings                                                                        |                   | ✅        | `progen` settings section                                                                                   |
|                                                                                 |                   |          |
| settings.http[<sup>**ⓘ**</sup>](#http_client)                                   |                   | ✅        | http client configuration                                                                                   |
| settings.http.debug                                                             |       bool        | ✅        | http client `DEBUG` mode                                                                                    |
| settings.http.base_url                                                          |      string       | ✅        | http client base `URL`                                                                                      |
| settings.http.headers                                                           | map[string]string | ✅        | http client base request `Headers`                                                                          |
| settings.http.query_params                                                      | map[string]string | ✅        | http client base request `Query Parameters`                                                                 |
|                                                                                 |                   |          |                                                                                                             |
| settings.groups[<sup>**ⓘ**</sup>](#groups_of_actions)                           |                   | ✅        | groups of actions                                                                                           |
| settings.groups.name                                                            |      string       | ✅        | group's name                                                                                                |
| settings.groups.actions                                                         |     []string      | ✅        | actions names                                                                                               |
| settings.groups.manual                                                          |       bool        | ✅        | determines that the group starts automatically (default `false`)                                            |
|                                                                                 |                   |          |                                                                                                             |
| dirs`<unique_suffix>`[<sup>**ⓘ**</sup>](#Generate)                              |     []string      | ✅        | list of directories to create                                                                               |
|                                                                                 |                   |          |                                                                                                             |
| rm`<unique_suffix>`[<sup>**ⓘ**</sup>](#rm)                                      |     []string      | ✅        | list for remove (files, dirs, all file in a dir)                                                            |
|                                                                                 |                   |          |                                                                                                             |
| <a name="files_actio_desk"><a/>files`<unique_suffix>`[<sup>**ⓘ**</sup>](#Files) |                   | ✅        | list file's `path` and `data`                                                                               |
| files.path                                                                      |      string       | ❌        | save file `path`                                                                                            |
| files.local                                                                     |      string       | `❕`      | local file path to copy                                                                                     |
| files.data                                                                      |      string       | `❕`      | save file `data`                                                                                            |
|                                                                                 |                   |          |                                                                                                             |
| files.get                                                                       |                   | `❕`      | struct describe `GET` request for getting file's data                                                       |
| files.get.url                                                                   |      string       | ❌        | request `URL`                                                                                               |
| files.get.headers                                                               | map[string]string | ✅        | request `Headers`                                                                                           |
| files.get.query_params                                                          | map[string]string | ✅        | request `Query Parameters`                                                                                  |
|                                                                                 |                   |          |                                                                                                             |
| cmd`<unique_suffix>`[<sup>**ⓘ**</sup>](#Commands)                               |                   | ✅        | configuration command list                                                                                  |
| cmd.exec                                                                        |      string       | ❌        | command to execution                                                                                        |
| cmd.args                                                                        |      []slice      | ✅        | list of command's arguments                                                                                 |
| cmd.dir                                                                         |      string       | ✅        | execution commands (`cmd.exec`) directory                                                                   |
|                                                                                 |                   |          |                                                                                                             |
| fs[<sup>**ⓘ**</sup>](#fs)                                                       |     []string      | ✅        | execute [text/template.Option](https://pkg.go.dev/text/template#Template.Option) on the list of directories |

`❕` only one must be specified in parent section

___

## Usage

### Generate

The cli executes commands and generate files and directories based on configuration file

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
2023-02-05 14:11:47	INFO	application working directory: /Users/user_1/GoProjects/service
2023-02-05 14:11:47	INFO	configuration file: progen.yml
2023-02-05 14:11:47	INFO	file process: x/some_file.txt
2023-02-05 14:11:47	INFO	dir created: x/y
2023-02-05 14:11:47	INFO	file saved: x/some_file.txt
2023-02-05 14:11:47	INFO	execute [dir: .]: touch second_file.txt
2023-02-05 14:11:47	INFO	execute [dir: .]: tree
out:
.
├── second_file.txt
└── x
    ├── some_file.txt
    └── y

2 directories, 2 files
```

### Execution

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

```console
% progen -v
2023-01-22 13:38:52	INFO	application working direcotry: /Users/user_1/GoProjects/service
2023-01-22 13:38:52	INFO	dir created: api/some_project/v1
2023-01-22 13:38:52	INFO	execute [dir: .]: chmod -R 777 api
2023-01-22 13:38:52	INFO	dir created: api/some_project_2/v1
2023-01-22 13:38:52	INFO	execute [dir: .]: chmod -R 777 api
```

### Templates

Configuration preprocessing uses [text/template](https://pkg.go.dev/text/template) of golang's `stdlib`.
Using templates could be useful to avoiding duplication in configuration file.
All `text/template` variables must be declared as comments and can be used only to configure data of configuration
file (all ones skipping for `file.data` section).
Configuration's `yaml` tag tree also use as `text/template` variables dictionary and can be use for avoiding duplication
in configuration file and files contents (`files` section).

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
    data: |
      {{$project_name}}

cmd:
  - "cat internal/{{$project_name}}.txt"
  - exec: ls
    dir: .
    args: [ -l ]
  - exec: tree
```

```console
% progen -v
2023-01-22 13:03:58	INFO	current working direcotry: /Users/user_1/GoProjects/service
2023-02-05 14:47:25	INFO	application working directory: /Users/user_1/GoProjects/service
2023-02-05 14:47:25	INFO	configuration file: progen.yaml
2023-02-05 14:47:25	INFO	file process: internal/SOME_PROJECT.txt
2023-02-05 14:47:25	INFO	file process: pkg/SOME_PROJECT-data/some_file.txt
2023-02-05 14:47:25	INFO	dir created: api/SOME_PROJECT/v1
2023-02-05 14:47:25	INFO	dir created: internal/some/file/path
2023-02-05 14:47:25	INFO	dir created: pkg/SOME_PROJECT-data
2023-02-05 14:47:25	INFO	file saved: internal/SOME_PROJECT.txt
2023-02-05 14:47:25	INFO	file saved: pkg/SOME_PROJECT-data/some_file.txt
2023-02-05 14:47:25	INFO	execute [dir: .]: cat internal/SOME_PROJECT.txt
out:
Project name:SOME_PROJECT

2023-02-05 14:47:25	INFO	execute [dir: .]: ls -l
out:
total 0
drwxr-xr-x  3 19798572  646495703   96 Feb  5 14:47 api
drwxr-xr-x  4 19798572  646495703  128 Feb  5 14:47 internal
drwxr-xr-x  3 19798572  646495703   96 Feb  5 14:47 pkg

2023-02-05 14:47:25	INFO	execute [dir: .]: tree
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

9 directories, 2 files
```

#### Custom template functions

| Function          |             args             | Description                                                                                                                                                                       |
|:------------------|:----------------------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `random`          |
| `random.Alpha`    |         length `int`         | Generates a random alphabetical `(A-Z, a-z)` string of a desired length.                                                                                                          | 
| `random.Num`      |         length `int`         | Generates a random numeric `(0-9)` string of a desired length.                                                                                                                    | 
| `random.AlphaNum` |         length `int`         | Generates a random alphanumeric `(0-9, A-Z, a-z)` string of a desired length.                                                                                                     |
| `random.ASCII`    |         length `int`         | Generates a random string of a desired length, containing the set of printable characters from the 7-bit ASCII set. This includes space (’ ‘), but no other whitespace character. |
| `slice`           |                              |                                                                                                                                                                                   |
| `slice.New`       |       N `any` elements       | Create new slice from any numbers of elements <br/>(`{ $element := slice.New "a" 1 "b" }}`)                                                                                       |
| `slice.Append`    | slice,<br/> N `any` elements | Add element to exists slice <br/>(`{{ $element := slice.Append $element "b"}}`)                                                                                                   |
| `strings`         |                              |                                                                                                                                                                                   |
| `strings.Replace` |  s, old, new string, n int   | Replace returns a copy of the string `s` with `old` replaced by `new` (work the same as `strings.Replace` from `stdlib`).                                                         |

Custom template's functions added as custom arguments to the template
[function map](https://pkg.go.dev/text/template#hdr-Functions).

---

## Flags

### <a name="config_file"><a/>Configuration file

By default `progen` try to find `progen.yml` file for execution. `-f` flag specify custom configuration file location:

```console
progen -f custom_conf.yaml
```

Instead of specifying a config file, you can pass a single configuration file in the pipe the file in via `STDIN`.
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

### <a name="print_err_trace"><a/>Print error stack trace

To print a stack trace of the error which occurred during execution of the `cli`,
use `-errtrace` flag:

```console
% progen -f ../not_exists_config.yml
2023-03-04 15:05:54	FATAL	read config: config file:
    github.com/kozmod/progen/internal/config.(*Reader).Read
        /Users/some_user/projects/progen/internal/config/reader.go:39
  - open ../not_exists_config.yml: no such file or directory
```

### <a name="dry_run"><a/>Dry Run mode

The `-dr` flag uses to execute configuration in dry run mod. All `action` will be executed without applying.

```yaml
## progen.yml

# {{$project_name := "SOME_PROJECT"}}
dirs:
  - api/{{ $project_name }}/v1 # apply template variables, but not create directories on 'dry run'
cmd:
  - tree # not execute on 'dry run' mode

files:
  - path: api/v1/some_file.txt # apply template variables and only printing the file's data
    data: |
      some file data data fot project: {{ $project_name }}
```

```console
% progen -v -dr
2023-03-07 07:57:52	INFO	application working directory: /Users/user_1/GoProjects/service
2023-03-07 07:57:52	INFO	configuration file: progen.yml
2023-03-07 07:57:52	INFO	file process: api/v1/some_file.txt
2023-03-07 07:57:52	INFO	dir created: api/SOME_PROJECT/v1
2023-03-07 07:57:52	INFO	execute [dir: .]: tree
2023-03-07 07:57:52	INFO	save file: create dir [api/v1] to store file [%!s(func() string=0x136ecc0)]
2023-03-07 07:57:52	INFO	file saved [path: api/v1/some_file.txt]:
some file data data fot project: SOME_PROJECT
2023-03-07 07:57:52	INFO	execution time: 3.69506ms
```

### <a name="awd"><a/>Application working directory

The `-awd` flag uses for setting application working directory.
All `paths` declared in the config file are calculated considering the root directory.

### <a name="print_config"><a/>Print configuration file

To print the configuration file after processing as [text/template](https://pkg.go.dev/text/template),
use `-printconf` flag:

```yaml
## progen.yml

vars:
  some_data: VARS_SOME_DATA

# {{- $var_1 := random.AlphaNum 15}}
#  {{- $var_2 := "echo some_%s"}}
cmd:
  - echo {{ $var_1 }}
  - "{{ printf $var_2  `value` }}"
  - echo {{ .vars.some_data }}
```

```console
% progen -printconf
2023-03-04 14:57:43	INFO	preprocessed config:
vars:
  some_data: VARS_SOME_DATA

#
#
cmd:
  - echo AHNsgyzVxRqeqLt
  - "echo some_value"
  - echo VARS_SOME_DATA
```

### <a name="files_preprocessing"><a/>Files preprocessing

By default, all files loading to the memory and process as [text/template](https://pkg.go.dev/text/template) before
saving to a file system.
To change this behavior, set `-pf=false`.

```console
% progen -v -dr -f progen.yml
2023-02-05 14:15:54	INFO	application working directory: /Users/user_1/GoProjects/service
2023-02-05 14:15:54	INFO	configuration file: progen.yml
2023-02-05 14:15:54	INFO	file process: api/v1/some_file.txt
2023-02-05 14:15:54	INFO	dir created: api/SOME_PROJECT/v1
2023-02-05 14:15:54	INFO	execute cmd: chmod -R 777 api/v1
2023-02-05 14:15:54	INFO	save file: create dir [api/v1] to store file [some_file.txt]
2023-02-05 14:15:54	INFO	file saved [path: api/v1/some_file.txt]:
some file data data fot project: SOME_PROJECT
```

### <a name="tvar"><a/>Template variables

Any part of template variable tree can be overrides using `-tvar` flag

```yaml
## progen.yml

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
2023-02-05 14:51:38	INFO	application working directory: /Users/user_1/GoProjects/service
2023-02-05 14:51:38	INFO	configuration file: progen.yml
2023-02-05 14:51:38	INFO	dir created: api/SOME_PROJECT/v1
2023-02-05 14:51:38	INFO	dir created: internal/some/file/path
2023-02-05 14:51:38	INFO	dir created: internal/overrided_path
```

### <a name="skip_actions"><a/>Skip `actions`

Set `-skip` flag to skip any `action` (only root actions: `cmd`, `files`, `dirs`). Value of the flag is a regular
expression.

```yaml
## progen.yml

dirs:
  - api/v1
cmd:
  - chmod -R 777 api/v1
dirs1:
  - api/v2
cmd1:
  - chmod -R 777 api/v2
dirs2:
  - api/v3
cmd2:
  - chmod -R 777 api/v3 
```

```console
% progen -v -dr -f progen.yml -skip=^dirs$ -skip=cmd.+ 
2023-02-05 14:18:11	INFO	application working directory: /Users/user_1/GoProjects/service
2023-02-05 14:18:11	INFO	configuration file: progen.yml
2023-02-05 14:18:11	INFO	action will be skipped: [cmd1]
2023-02-05 14:18:11	INFO	action will be skipped: [cmd2]
2023-02-05 14:18:11	INFO	action will be skipped: [dirs]
2023-02-05 14:18:11	INFO	execute cmd: chmod -R 777 api/v1
2023-02-05 14:18:11	INFO	dir created: api/v2
2023-02-05 14:18:11	INFO	dir created: api/v3
```

---

## Actions and tags

### <a name="http_client"></a>Http Client

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

### <a name="groups_of_actions"></a>Groups of actions

All actions execute in declaration order in the config file and can be union to groups.
All actions in `manual` groups will be skipped during execution process.

```yaml
settings:
  groups:
    - name: group1
      actions: [ cmd, cmd_2 ]
      manual: true
    - name: group2
      actions: [ cmd_2 ]
      manual: true

cmd:
  - echo CMD_1

cmd_2:
  - echo CMD_2

cmd_3:
  - echo CMD_3

cmd_4:
  - echo CMD_4

```

```console
% progen -v
2024-02-05 23:08:21     INFO    application working directory: /Users/user_1/GoProjects/service
2024-02-05 23:08:21     INFO    configuration file: progen.yml
2024-02-05 23:08:21     INFO    manual actions will be skipped: [cmd, cmd_2]
2024-02-05 23:08:21     INFO    execute [dir: .]: echo CMD_3
out:
CMD_3

2024-02-05 23:08:21     INFO    execute [dir: .]: echo CMD_4
out:
CMD_4

2024-02-05 23:08:21     INFO    execution time: 7.916615ms
```

Actions in `manual` groups execute using `gp` flag (all action execute only once independent on declaration's quantity
in different groups).

```console
% progen -v -gp=group1 -gp=group2
2024-02-05 23:19:50     INFO    application working directory: /Users/user_1/GoProjects/service
2024-02-05 23:19:50     INFO    configuration file: progen.yml
2024-02-05 23:19:50     INFO    groups will be execute: [group1, group2]
2024-02-05 23:19:50     INFO    execute [dir: .]: echo CMD_1
out:
CMD_1

2024-02-05 23:19:50     INFO    execute [dir: .]: echo CMD_2
out:
CMD_2

2024-02-05 23:19:50     INFO    execution time: 7.192257ms

```

### Files

File's content can be declared in configuration file (`files.data` tag) or
can be received from local (`files.local`) or remote (`files.get`) storage.
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
    # GET file from remote storage (using common http client config)
    get:
      # reuse `base` URL of common http client config (http.base_url)
      url: Dockerfile/raw?ref=feature/project_templates"
```

```console
% progen -v
2023-02-05 14:47:25	INFO	current working direcotry: /Users/user_1/GoProjects/service
2023-02-05 14:47:25	INFO	configuration file: progen.yaml
2023-02-05 14:47:25	INFO	file process: files/Readme.md
2023-02-05 14:47:25	INFO	file process: files/.gitignore
2023-02-05 14:47:25	INFO	file process: files/.editorconfig
2023-02-05 14:47:25	INFO	file process: files/.gitlab-ci.yml
2023-02-05 14:47:25	INFO	file process: files/Dockerfile
...
2023-02-05 14:47:25	INFO	file saved: files/Readme.md
2023-02-05 14:47:25	INFO	file saved: files/.gitignore
2023-02-05 14:47:25	INFO	file saved: files/.editorconfig
2023-02-05 14:47:25	INFO	file saved: files/.gitlab-ci.yml
2023-02-05 14:47:25	INFO	file saved: files/Dockerfile
...
```

### Commands

Execution commands process configured by specifying __commands working directory__ and commands definition.
Default value of __commands working directory__ (`dir` tag) is `.`.
__Commands working directory__ calculate from the __application working directory__.

```yaml
## progen.yml

cmd:
  - exec: ls -l
    args: [ - l ]
    dir: .github/workflows
  - exec: tree
    args: [ -L, 1 ]
```

```console
% progen -v 
2023-02-02 22:18:20	INFO	application working directory: /Users/user_1/GoProjects/progen
2023-02-02 22:18:20	INFO	configuration read: progen.yml
2023-02-02 22:18:20	INFO	execute [dir: .github/workflows]: ls -l
out:
total 16
-rw-r--r--  1 19798572  646495703  762 Feb  1 09:15 release.yml
-rw-r--r--  1 19798572  646495703  377 Jan 24 20:06 test.yml

2023-02-02 22:18:20	INFO	execute [dir: .]: tree -L 1
out:
.
├── LICENSE
├── Makefile
├── Readme.md
├── go.mod
├── go.sum
├── internal
├── main.go
└── tmp

2 directories, 6 files
```

`cmd` action maintains "short" declaration syntax

```yaml
## progen.yml

cmd:
  - pwd
  - ls -a
```

```console
% progen -v -dr
2023-02-15 17:56:58	INFO	application working directory: /Users/user_1/GoProjects/progen
2023-02-15 17:56:58	INFO	configuration file: short.yml
2023-02-15 17:56:58	INFO	execute [dir: .]: pwd
2023-02-15 17:56:58	INFO	execute [dir: .]: ls -a
```

### <a name="fs"></a>File System

`fs` section configure execution [text/template](https://pkg.go.dev/text/template) on a directories tree.
All files in the `tree` processed as `template`. Files and directories names also could be configured as templates.

```yaml
## progen.yml

var_d: VAR_d
var_f: VAR_f

cmd:
  - cp -a ../asserts/. ../out/
  - exec: tree
    dir: .

fs:
  - test_dir
  - test_dir_2

cmd_finish:
  - exec: tree
    dir: .
```

```console
% progen -v -awd=out -f ../progen.yml
2023-02-12 14:01:45	INFO	application working directory: /Users/user_1/GoProjects/progen
2023-02-12 14:01:45	INFO	configuration file: ../progen.yml
2023-02-12 14:01:45	INFO	execute [dir: .]: cp -a ../asserts/. ../out/
2023-02-12 14:01:45	INFO	execute [dir: .]: tree
out:
.
├── test_dir
│   ├── file1
│   └── {{ .var_d }}
│       └── {{ .var_f }}
└── test_dir_2
    ├── file1
    └── {{ .var_d }}
        └── {{ .var_f }}

4 directories, 4 files

2023-02-12 14:01:45	INFO	dir created: test_dir/VAR_d
2023-02-12 14:01:45	INFO	file saved: test_dir/file1
2023-02-12 14:01:45	INFO	file saved: test_dir/VAR_d/VAR_f
2023-02-12 14:01:45	INFO	dir created: test_dir_2/VAR_d
2023-02-12 14:01:45	INFO	file saved: test_dir_2/file1
2023-02-12 14:01:45	INFO	file saved: test_dir_2/VAR_d/VAR_f
2023-02-12 14:01:45	INFO	fs: remove: test_dir_2/{{ .var_d }}/{{ .var_f }}
2023-02-12 14:01:45	INFO	fs: remove: test_dir_2/{{ .var_d }}
2023-02-12 14:01:45	INFO	execute [dir: .]: tree
out:
.
├── test_dir
│   ├── VAR_d
│   │   └── VAR_f
│   └── file1
└── test_dir_2
    ├── VAR_d
    │   └── VAR_f
    └── file1
```

### <a name="rm"></a>Rm

`rm` use to remove files, directories or files inside a directory.

```yaml
rm:
  # remove the dir
  - some_dir
  # remove all files in the dir
  - some_dir_2/*
  # remove the file
  - some_dir_3/file.txt
```

```console
% progen -v
2024-02-09 22:50:51     INFO    application working directory: /Users/user_1/GoProjects/progen
2023-02-12 14:01:45     INFO    configuration file: progen.yml
2024-02-09 22:50:51     INFO    rm: some_dir
2024-02-09 22:50:51     INFO    rm all: some_dir_2/*
2024-02-09 22:50:51     INFO    rm: some_dir_3/file.txt
2024-02-09 22:50:51     INFO    execution time: 350.149µs
```

---

### <a name="lib_usage"><a/>Lib

To use `progen` for building custom generator based on `go` language, imports `pkg/core` package
and implements required algorithm:
```golang
package main

import (
	"log"
	"testing/fstest"

	"github.com/kozmod/progen/pkg/core"
)

var fs = fstest.MapFS{
	"1": {
		Data: []byte("aaaa"),
	},
}

func main() {
	// Parse default config base on flags
	c, err := core.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	var e core.Engin

	// add actions
	e.AddActions(
		// create files actions
		core.FilesAction(
			"create_file",
			core.File{Path: "./xx/1", Data: []byte("file_1")},
			core.File{Path: "./xx/2", Data: []byte("file_2")},
			core.File{Path: "./xx/rm_1", Data: []byte("file_rm")},
		).WithPriority(1),
		// rm actions
		core.RmAction(
			"rm",
			"./xx/rm_1",
		).WithPriority(2),
	)
	e.AddActions(
		// create file system action
		core.FsSaveAction(
			"fs_1",
			core.TargetFs{
				TargetDir: "./xx/fs",
				Fs:        fs, // any [io/fs.FS] (from local system, embed, etc.)
			},
		).WithPriority(3),
		
		// cmd action
		core.CmdAction(
			"tree_1",
			core.Cmd{
				Cmd: "tree",
				Dir: "./xx",
			},
		).WithPriority(4),
	)

	err = e.Run(c)
	if err != nil {
		log.Fatal(err)
	}
}

```
```console
% go build .
% ./tmp -v  
2024-12-11 10:16:17     INFO    action is going to be execute ('priopiry':'name')['1':'create_file','2':'rm','3':'fs_1','4':'tree_1']
2024-12-11 10:16:17     INFO    file saved: xx/1
2024-12-11 10:16:17     INFO    file saved: xx/2
2024-12-11 10:16:17     INFO    file saved: xx/rm_1
2024-12-11 10:16:17     INFO    rm: ./xx/rm_1
2024-12-11 10:16:17     INFO    dir created: xx/fs
2024-12-11 10:16:17     INFO    file saved: xx/fs/1
2024-12-11 10:16:17     INFO    execute [dir: ./xx]: tree
out:
.
├── 1
├── 2
└── fs
    └── 1

1 directory, 3 files
```

---

### Examples

[progen-example](https://github.com/kozmod/progen-examples) repository contains useful examples of usage cli
