package cli_test

import (
	"github.com/silvertern/cli"
	"os"
	"testing"
)

func setup_usage_app() cli.App {
	co := cli.NewCommand("checkout", "Check out a branch or revision").
		WithShortcut("co").
		WithArg(cli.NewArg("revision", "branch or revision to checkout")).
		WithArg(cli.NewArg("fallback", "branch to fallback").AsOptional()).
		WithOption(cli.NewOption("branch", "create branch if missing").WithChar('b').WithType(cli.TypeBool)).
		WithAction(func(args []string, options map[string]string) int {
			return 25
		}).
		WithCommand(cli.NewCommand("sub-cmd1", "First sub-command")).
		WithCommand(cli.NewCommand("sub-cmd2", "Second sub-command"))

	rmt := cli.NewCommand("remote", "Work with git remotes")

	return cli.New("git tool").
		WithCommand(co).
		WithCommand(rmt).
		WithOption(cli.NewOption("verbose", "Verbose execution").WithChar('v').WithType(cli.TypeBool))
}

type stringwriter struct {
	str string
}

func (s *stringwriter) Write(p []byte) (n int, err error) {
	s.str = s.str + string(p)
	return len(p), nil
}

func TestApp_Usage_NestedCommandHelp_ok(t *testing.T) {
	a := setup_usage_app()
	w := &stringwriter{}
	a.Run([]string{"./foo", "co", "-hb", "5.5.5"}, w)
	expected := `foo checkout [--verbose] [--branch] <revision> [fallback]

Description:
    Check out a branch or revision

Arguments:
    revision                branch or revision to checkout
    fallback                branch to fallback, optional

Options:
    -v, --verbose           Verbose execution
    -b, --branch            create branch if missing

Sub-commands:
    foo checkout sub-cmd1   First sub-command
    foo checkout sub-cmd2   Second sub-command
`
	assertAppUsageOk(t, expected, w.str)
}

func TestApp_Usage_NestedCommandParsginError_ok(t *testing.T) {
	a := setup_usage_app()
	w := &stringwriter{}
	a.Run([]string{"./foo", "co"}, w)
	expected := `fatal: missing required argument revision
usage: foo checkout [--verbose] [--branch] <revision> [fallback]
`
	assertAppUsageOk(t, expected, w.str)
}

func TestApp_Usage_TopWithNoAction(t *testing.T) {
	a := setup_usage_app()
	w := &stringwriter{}
	code := a.Run([]string{"./foo"}, w)
	if code != 1 {
		t.Errorf("expected exit code 1, found %v", code)
	}
	expected := `foo [--verbose]

Description:
    git tool

Options:
    -v, --verbose   Verbose execution

Sub-commands:
    foo checkout    Check out a branch or revision
    foo remote      Work with git remotes
`
	assertAppUsageOk(t, expected, w.str)
}

func TestApp_Usage_README(t *testing.T) {
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
		WithArg(cli.NewArg("remote", "remote to add"))

	rmt := cli.NewCommand("remote", "Work with git remotes").
		WithCommand(add)

	app := cli.New("git tool").
		WithOption(cli.NewOption("verbose", "Verbose execution").WithChar('v').WithType(cli.TypeBool)).
		WithCommand(co).
		WithCommand(rmt)
		// no action attached, just print usage when executed

	app.Run([]string{"./gitc", "co", "-f", "dev"}, os.Stdout)
}

func assertAppUsageOk(t *testing.T, expectedOutput, actualOutput string) {
	if expectedOutput != actualOutput {
		t.Errorf("expected output: %v, found: %v", expectedOutput, actualOutput)
	}
}
