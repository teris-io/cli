[![Build status][buildimage]][build] [![Coverage][codecovimage]][codecov] [![GoReportCard][cardimage]][card] [![API documentation][docsimage]][docs]

# cli

Module `cli` provides a simple, fast and complete API for building command line applications in Go.
In contrast to other libraries additional emphasis is put on the definition and validation of
positional arguments and consistent usage outputs combining options from all command levels into
one block.

```
	co := cli.NewCommand("checkout", "checkout a branch or revision").
		WithShortcut("co").
		WithArg(cli.NewArg("branch")).
		WithOption(cli.NewOption("branch", "Create branch").WithChar('b').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("upstream", "Set upstream").WithChar('u').WithType(cli.TypeBool)).

	add := cli.NewCommand("add", "add a remote").
		WithArg(cli.NewArg("remote")).

	rmt := cli.NewCommand("remote", "operations with remotes").WithCommand(add)

	return cli.New("git tool").
		WithCommand(co).
		WithCommand(rmt)
```

[docs]: https://godoc.org/github.com/silvertern/cli
[docsimage]: http://img.shields.io/badge/godoc-reference-blue.svg?style=flat

[build]: https://travis-ci.org/silvertern/cli
[buildimage]: https://travis-ci.org/silvertern/cli.svg?branch=master

[codecov]: https://codecov.io/gh/silvertern/cli
[codecovimage]: https://codecov.io/gh/silvertern/cli/branch/master/graph/badge.svg

[card]: https://goreportcard.com/report/github.com/silvertern/cli
[cardimage]: https://goreportcard.com/badge/github.com/silvertern/cli
