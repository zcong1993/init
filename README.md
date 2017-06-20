# init
[![Go Report Card](https://goreportcard.com/badge/github.com/zcong1993/init)](https://goreportcard.com/report/github.com/zcong1993/init)
[![Build Status](https://travis-ci.org/zcong1993/init.svg?branch=master)](https://travis-ci.org/zcong1993/init)

> Init project from git repo template

## Install

Download the latest [release](https://github.com/zcong1993/init/releases), then place to your `$PATH` folder and rename it to `init`.

**Note:** *nix user maybe need `chmod +x init` to make it executable.

## Usage

```bash
$ init [options] gituser/repo ./your/local/folder
# options
# -install, -i Download template from github, replace cache if exists
# -force, -f   Remove local folder if exists
# example
$ init zcong1993/template-go ./gotest
```

## Create a template

### 1. create a github repo
should require `template` and `init.json` in root folder, like this:
```
|-template
    |-tpl.go
    |-README.md
|-init.json
```
### 2. custom a init.json file, like this:
```json
{
  "prompts": {
    "name": {
      "message": "Your project name ?",
      "type": "string"
    },
    "description": {
      "message": "How would you descripe the new project ?",
      "type": "string",
      "default": "my go project"
    },
    "username": {
      "message": "Your github username ?",
      "type": "string"
    },
    "cli": {
      "message": "Is a cli project ?",
      "type": "confirm"
    },
    "test": {
      "message": "Choose test :",
      "type": "list",
      "choices": [
        "travis",
        "wercker",
        "none"
      ]
    }
  },
  "filters": {
    "main_test.go": "test != \"none\"",
    ".travis.yml": "test == \"travis\"",
    "wercker.yml": "test == \"wercker\"",
    "build.sh": "cli",
    "Makefile": "cli"
  }
}
```
`prompts` provide `data` for template files by inquiring. And type can be `string, confirm and list`. `string` can have a `default` value, if not have will be a required data. `list` should have a `choices` list.

`filters` can controll which file should be used, only when right expression is `true` left file can be generating. Left file use [doublestar](https://github.com/bmatcuk/doublestar) so it support `some/path/*` and `some/**`. The right expression can use all values `prompts` provided cause it will be eval after inquire.

*full example* please see [zcong1993/template-go](https://github.com/zcong1993/template-go)

### 3. put your template files in repo template folder
all file in template folder will be compiled as `text/template` does, so you can use all `golang template` can use except custom function.

*full example* please see [zcong1993/template-go](https://github.com/zcong1993/template-go)

### 4. publish to github

## License

MIT &copy; zcong1993
