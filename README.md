# testecho [![GoDoc](https://godoc.org/github.com/codeactual/testecho?status.svg)](https://pkg.go.dev/mod/github.com/codeactual/testecho) [![Go Report Card](https://goreportcard.com/badge/github.com/codeactual/testecho)](https://goreportcard.com/report/github.com/codeactual/testecho) [![Build Status](https://travis-ci.org/codeactual/testecho.png)](https://travis-ci.org/codeactual/testecho)

testecho is a program to assist test cases which assert the subject starts and ends processes as expected. Its flags allow selection of stdout, stderr, exit code, etc.

The [testecho package](https://pkg.go.dev/mod/github.com/codeactual/testecho) provides functions to more easily run the CLI from test cases.

# Usage

> To install: `go get -v github.com/codeactual/testecho/cmd/testecho`

## Examples

> Display help

```bash
testecho --help
```

> Print "out" to standard output:

```bash
testecho --stdout out
```

> Same as above but also print "err" to standard error:

```bash
testecho --stdout out --stderr err
```

> Same as above but also exit with code 7 instead of 0:

```bash
testecho --stdout out --stderr err --code 7
```

> Same as above but also sleep for 5 seconds after printing:

```bash
testecho --stdout out --stderr err --code 7 --sleep 5
```

> Spawn another testecho proecss, print its PID, and then sleep "forever" (10000 seconds):

```bash
testecho --spawn
```

> Same as above but also print "err" to standard error:

```bash
testecho --spawn --stderr err
```

> Print standard input:

```bash
echo "out" | testecho
```

> Same as above but also print "err" to standard error:

```bash
echo "out" | testecho --stderr err
```

> Same as above but also exit with code 7 instead of 0:

```bash
echo "out" | testecho --stderr err --code 7
```

> Same as above but also sleep for 5 seconds after printing

```bash
echo "out" | testecho --stderr err --code 7 --sleep 5
```

## Examples (in other projects)

- [codeactual/boone](https://sourcegraph.com/search?q=repo:%5Egithub%5C.com/codeactual/boone%24+testecho%7Cechopath)

# Development

## License

[Mozilla Public License Version 2.0](https://www.mozilla.org/en-US/MPL/2.0/) ([About](https://www.mozilla.org/en-US/MPL/), [FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/))

## Contributing

- Please feel free to submit issues, PRs, questions, and feedback.
- Although this repository consists of snapshots extracted from a private monorepo using [transplant](https://github.com/codeactual/transplant), PRs are welcome. Standard GitHub workflows are still used.
