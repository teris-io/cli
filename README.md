[![Build status][buildimage]][build] [![Coverage][codecovimage]][codecov] [![GoReportCard][cardimage]][card] [![API documentation][docsimage]][docs]

# Simple and complete API for building command line applications in Go

Module `cli` provides a simple, fast and complete API for building command line applications in Go.
In contrast to other libraries the emphasis is put on the definition and validation of
positional arguments, handling of options from all levels in a single block as well as
a minimalistic set of dependencies.

The core of the module is the command, option and argument parsing logic. After a successful parsing the 
command action is evaluated passing a slice of (validated) positional arguments and a map of (validated) options.
No more no less.

## Definition

```
co := cli.NewCommand("checkout", "checkout a branch or revision").
  WithShortcut("co").
  WithArg(cli.NewArg("revision", "branch or revision to checkout")).
  WithOption(cli.NewOption("branch", "Create branch if missing").WithChar('b').WithType(cli.TypeBool)).
  WithOption(cli.NewOption("upstream", "Set upstream for the branch").WithChar('u').WithType(cli.TypeBool)).
  WithAction(func(args []string, options map[string]string) int {
    // do something
    return 0
  })

add := cli.NewCommand("add", "add a remote").
  WithArg(cli.NewArg("remote", "remote to add")).

rmt := cli.NewCommand("remote", "Work with git remotes").
  WithCommand(add)

app := cli.New("git tool").
  WithOption(cli.NewOption("verbose", "Verbose execution").WithChar('v').WithType(cli.TypeBool)).
  WithCommand(co).
  WithCommand(rmt)
  // no action attached, just print usage when executed

os.Exit(app.Run(os.Args, os.Stdout))
```

## Execution

Given the above definition is for a git client, e.g. `gitc`, running `gitc` with no arguments or with `-h` will
produce (the exit code will be 1 in the former case, because the action is missing, and 0 in the latter, because
help explicitly requested):

```
gitc [--verbose]

Description:
    git tool

Options:
    -v, --verbose   Verbose execution

Sub-commands:
    git checkout    checkout a branch or revision
    git remote      Work with git remotes
```

Running `gitc` with arguments matching e.g. the `checkout` definition, `gitc co -vbu dev` or
`gitc checkout -v --branch -u dev` will execute the command as expected. Running into a parsing error, e.g.
 by providing an unknown option `gitc co -f dev`, will output a parsing error and a short usage string:

```
fatal: unknown flag -f
usage: gitc checkout [--verbose] [--branch] [--upstream] <revision>
```


### License and copyright

	Copyright (c) 2017. Oleg Sklyar and teris.io. MIT license applies. All rights reserved.


[build]: https://travis-ci.org/teris-io/cli
[buildimage]: https://travis-ci.org/teris-io/cli.svg?branch=master

[codecov]: https://codecov.io/github/teris-io/cli?branch=master
[codecovimage]: https://codecov.io/github/teris-io/cli/coverage.svg?branch=master

[card]: http://goreportcard.com/report/teris-io/cli
[cardimage]: https://goreportcard.com/badge/github.com/teris-io/cli

[docs]: https://godoc.org/github.com/teris-io/cli
[docsimage]: http://img.shields.io/badge/godoc-reference-blue.svg?style=flat
