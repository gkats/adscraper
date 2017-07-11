# adscripted scraper

All you need to scrape raw ads.

## Setup

#### 1. Install Go 1.8.x

Go to [golang.org downloads page](https://golang.org/dl/), download the binary for your architecture and then follow the [installation instructions](https://golang.org/doc/install). It's probably as simple as (running with super user privileges)
```
$ tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```
It's recommended to use the default installation location and __not__ set the `GOROOT` environment variable.

#### 2. Add Go's install location to your `PATH`.

Add this to your shell's init scripts.
```
$ export PATH=$PATH:/usr/local/go/bin
```

#### 3. Install PostgreSQL 9.6

Go to the official PostgreSQL [downloads](https://www.postgresql.org/download/) page. Follow the instructions for your operating system. If you're on Linux, PostgreSQL will most likely be already included in your package manager sources. If you're on macOS you can use [Homebrew](https://brew.sh/).

Don't forget to install both the server and the client libraries.

#### 4. Clone the project

```
$ git clone git@github.com:gkats/keywords.git
```

#### 5. Set the `GOPATH` environment variable

You can set the `GOPATH` environment variable to whatever you like, however the project assumes that it's cloned under `$(GOPATH)/src/github.com/gkats/`. There is a helper file provided to automatically set the `GOPATH` variable.

Go to the root directory and source the `.gopath` file.
```
$ . .gopath
```
This sets the `GOPATH` environment variable to the project root.

Happy hacking!

## Contribute

The package contains a `Makefile` for building with GNU Make. There are various targets in the Makefile. The default just builds (compiles) the package.

Before you run any Makefile rules, you need to set your `GOPATH`. The `GOPATH` variable's value is tightly related to the way you've set up your project directory hierarchy. It is recommended that you follow the `$(GOPATH)/src/github.com/gkats/scraper` directory structure. All you have to do then is `source .gopath` from the project root. While you can source the file only once, we'll include the directive for every run of make.

1. To install the package run
```
$ . .gopath && make install
# or
$ GOPATH=/path/to/gopath make install
```
The above command also runs `gofmt` and `govet` before installing the package.

2. To build the package run
```
$ . .gopath && make build
# or
$ GOPATH=/path/to/gopath make build
```

3. You can also run `gofmt` and `govet`
```
$ . .gopath && make fmt
# or
$ GOPATH=/path/to/gopath make fmt
```

```
$ . .gopath && make vet
# or
$ GOPATH=/path/to/gopath make vet
```

## Test

Nothing here at this point.

## Run

There are three separate programs bundled in the repo.


__keywords__
The main keywords program reads keywords from a file and stores them into the database. Once installed, you can invoke the program with
```
$ $(GOPATH)/bin/keywords -f absolute/path/to/keywords/file -d user:password\@host:port/database
```

You need to create a keywords file first. For an example see the sample `./keywords.dat.sample`.

When you're in doubt just run `$ $(GOPATH)/bin/keywords --help`.

__server__
This is an HTTP server used to read keywords, store ads and update the keywords scraping data. Run the program with
```
$ $(GOPATH)/bin/server -d user:password\@host:port/database
```
Run `$ $(GOPATH)/bin/server --help` for more information.

__scraper__
The application that scrapes raw ads from google results. It performs a request to get least scraped keywords (random), queries google for results and then posts them back to the server. Run it with
```
$ $(GOPATH)/bin/scraper -h https://server.hostname
```
Run `$ $(GOPATH)/bin/scraper --help` for more information.


