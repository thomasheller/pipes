# pipes

Pipes provides support for long/complex command pipelines.

## Setup

### Prerequisites

Install Go for your platform: https://golang.org

I recommend gimme: https://github.com/travis-ci/gimme

```
mkdir -p ~/bin
curl -sL -o ~/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/bin/gimme
```

### Download

Clone this repository into the `src` directory of your
[`$GOPATH`](https://golang.org/doc/code.html#GOPATH)
(default: `~/go/src`):

```
git clone https://github.com/thomasheller/pipes ~/go/src/github.com/thomasheller/pipes
```

### Build

```
eval $(~/bin/gimme stable)
cd ~/go/src/github.com/thomasheller/pipes
go get
go build
go install
```

### Local setup

The binary `pipes` will appear in the `bin` directory of your `$GOPATH`
(default: `~/go/bin`).

You can either move the `pipes` binary to a directory that is in your `$PATH`
or add the Go `bin` directory to your `$PATH`.

## Usage

### Examples

```
cd examples

pipes -from example1.txt -pipe reverse.pipe

pipes -from example1.txt -pipe replace.pipe

pipes -from example2.txt -pipe foo.pipe
```

## Improvements (ToDo)

- Output intermediate pipe stages for debugging.
- Cache pipe stages and re-use intermediate stages when only parts of the pipe are changed.

