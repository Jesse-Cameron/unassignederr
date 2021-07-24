# `unassignederrcheck`

`unassignederrcheck` is a tool for checking for returned errors which have been left unassigned or uninitialized.

Golang draws a distinction betweens a between a nil and a nil interface. This behaviour has been well documented in the [Golang FAQ](https://golang.org/doc/faq#nil_error). It's something that easily could be missed. This linting rule is indended to highlight cases where you may accidentally return an error struct that hasn't been initialized.

**Bad 😿**

```golang
type NaughtyErr struct {
      Msg string
}

func (e *NaughtyErr) Error() string { return e.Msg }

func returnsError() error {
      var err *NaughtyErr

      return err // err is unassigned
}
```

**Better 😸**

```golang
func returnsError() error {
      var err *NaughtyErr

      if doSomething() == false {
            err = &NaughtyErr{Msg: "error message"}
      }

      return err // err is correctly assigned
}
```

## Install

```bash
go get -u github.com/Jesse-Cameron/golang-nil-error-struct
```

## Usage

```
nil_error_struct: A tool for identifying when a uninitialised error struct is being incorrectly returned.

Usage: nil_error_struct [-flag] [package]


Flags:
  -V    print version and exit
  -all
        no effect (deprecated)
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -fix
        apply all suggested fixes
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -source
        no effect (deprecated)
  -tags string
        no effect (deprecated)
  -trace string
        write trace log to this file
  -v    no effect (deprecated)
```


**Examples**
```
$ nil-err-check ./...
```

```
$ nil-err-check github.com/package_name
```

**go/analysis**

The package provides Analyzer instance that can be used with go/analysis API.
