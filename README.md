# gocopy

## What is it?

`gocopy` is a simple utility to copy functions and types from Go source files.

## Why?

Because sometimes I want to copy over some code from a file somewhere else.
It's easy to do the first copy & paste, but maintenance becomes a burden as
the source gets updated.

So I've built this tool to be used with `go generate`.

## Examples

### Copying a type

To copy the `Decoder` type from the `encoding/json` package:

`gocopy type -u "https://golang.org/src/encoding/json/decode.go?m=text" -n Number`

If you want to include all its methods, use the `-m` flag (shorthand for `--methods`):

`gocopy type -u "https://golang.org/src/encoding/json/decode.go?m=text" -n Number -m`

### Copying a function

To copy the `Decode` function from the `encoding/json` package:

`gocopy function -u "https://golang.org/src/encoding/json/decode.go?m=text" -n Unmarshal`