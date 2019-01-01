# Gopher
[![Go Report Card](https://goreportcard.com/badge/github.com/MovieStoreGuy/gopher)](https://goreportcard.com/report/github.com/MovieStoreGuy/gopher)
[![Maintainability](https://api.codeclimate.com/v1/badges/3f34c6090b7f4fa04596/maintainability)](https://codeclimate.com/github/MovieStoreGuy/gopher/maintainability)
[![Build Status](https://travis-ci.org/MovieStoreGuy/gopher.svg?branch=master)](https://travis-ci.org/MovieStoreGuy/gopher)
[![codecov](https://codecov.io/gh/MovieStoreGuy/gopher/branch/master/graph/badge.svg)](https://codecov.io/gh/MovieStoreGuy/gopher)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)  
_Gopher enables the management and development of Golang projects while switch between different developer environments._

When trying to develop projects under different usernames or organisations, it can become frustrating walking through the directories. In order to simplify creating projects, gopher aims to make creating a project less taxing on the developer creating it and just let them get on with it.

## Installation
```
# With out using go modules
mkdir -p ${GOPATH:-'~/go'}/src/github.com/MovieStoreGuy
git clone https://github.com/MovieStoreGuy/gopher.git ${GOPATH:-'~/go'}/src/github.com/MovieStoreGuy/gopher
cd ${GOPATH:-'~/go'}/src/github.com/MovieStoreGuy/gopher
dep ensure
go install

# Using go modules
git clone https://github.com/MovieStoreGuy/gopher.git
cd gopher
go install
```

## Usage
To get started with `gopher` make sure it is on your $PATH. Once it is on the path, you can start using it.
```
Usage: gopher [global options] <verb> [verb options] [nested verb...]

Global options:
        -p, --profile      Define the profile to use with gopher (default: current)
        -v, --verbose      Enable verbose logging
        -h, --help         Show the global help

Verbs:
    create:
        -p, --path         Define an alternate path to create the project (Requires Go v1.11+)
        -h, --help         creates a project based on the supplied profile with the name added after flags
    profile:
        -n, --name         The name of the profile to store
        -v, --vcs          The name of the VCS to use
        -u, --username     The username to develop as
        -m, --make-default Make the defined profile the default used
        -h, --help         To use profile ensure you supply a subverb of [create, show, set]
    project:
        -h, --help         To use projects ensure you supply a subverb of [show]
```
**Note: Due to a limitation with the library being used within, options must occur before the verb.**

## Further Development
This project is in alpha and should be used with care. Any feature requests should be raised against the project.
