# gocd
__Change directory to a Go package__

A very simple command line application to automatically change directory based on a Go package name
Inpirated by github.com/relvacode/gocd

## Install

```bash
$ go get -v github.com/Ak-Army/go-cd
$ cat `go env GOPATH`/src/github.com/Ak-Army/go-cd/bashrc >> ~/.bashrc
```

  * Run `go get -v github.com/Ak-Army/go-cd` to install package dependencies
  * Add the contents of `bashrc` to your `~/.bashrc`
  * Either `source ~/.bashrc` or re-open your terminal window


## Usage

#### Absolute Package Names

You can navigate to a Go package directly

```bash
$ gocd github.com/Ak-Army/go-cd
```

#### Fuzzy Package Names

You can also use a fuzzy match for the package you want

```bash
$ gocd username/pkg
$ gocd pkg
```

gocd will scan your `GOPATH` and look for matches, if one match is found then you are taken to it. 

If more than one match is found supply the requested index as the second argument.

```bash
$ gocd txt
  0 golang.org/x/text
  1 golang.org/x/text/cases
  2 golang.org/x/text/cmd/gotext
  3 golang.org/x/text/cmd/gotext/examples/extract
  4 golang.org/x/text/cmd/gotext/examples/extract_http
  5 golang.org/x/text/cmd/gotext/examples/extract_http/pkg
  
$ gocd txt 0
```

#### Change Directory to Vendor Parent

```bash
$ gocd ^
```

Using `^` will navigate to the parent package of a vendored directory

##### GOPATH

Go to the `GOPATH` by calling gocd without arguments

```bash
$ gocd
```


## Develop
```bash
$ go get github.com/Ak-Army/go-cd
$ make init
```

  * Run `go get github.com/Ak-Army/go-cd` to install package dependencies
  * Run `make init` for initialize dependecies
  * Run `make full-test` for test
  * Run `make build` for build
