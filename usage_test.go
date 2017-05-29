package cli_test

import (
	"github.com/silvertern/cli"
	"testing"
)

func setup_usage_app() cli.App {
	co := cli.NewCommand("checkout", "Check out a branch or revision").
		WithShortcut("co").
		WithArg(cli.NewArg("revision", "branch or revision to checkout")).
		WithOption(cli.NewOption("branch", "create branch if missing").WithChar('b').WithType(cli.TypeBool)).
		WithAction(func(args []string, options map[string]string) int {
			return 25
		}).
		WithCommand(cli.NewCommand("dummy1", "First dummy command")).
		WithCommand(cli.NewCommand("dummy2", "Second dummy command"))

	add := cli.NewCommand("add", "add a remote").
		WithOption(cli.NewOption("force", "Force").WithChar('f').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("quiet", "Quiet").WithChar('q').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("default", "Default"))

	rmt := cli.NewCommand("remote", "operations with remotes").WithCommand(add)

	return cli.New("git tool").
		WithCommand(co).
		WithCommand(rmt)
}

type stringwriter struct {
	str string
}

func (s *stringwriter) Write(p []byte) (n int, err error) {
	s.str = s.str + string(p)
	return len(p), nil
}

func TestApp_Usage_NestedCommandHelp_ok(t *testing.T) {
	a := setup_run_app()
	w := &stringwriter{}
	a.Run([]string{"./foo", "co", "-hb", "5.5.5"}, w)
	expected := `foo checkout [--branch] <branch>

Description:
    Check out a branch or revision

Arguments:
    branch                branch to checkout

Options:
    -b, --branch          Create branch if missing

Sub-commands:
    foo checkout dummy1   First dummy command
    foo checkout dummy2   Second dummy command
`
	assertAppUsageOk(t, expected, w.str)
}

func TestApp_Usage_NestedCommandParsginError_ok(t *testing.T) {
	a := setup_run_app()
	w := &stringwriter{}
	a.Run([]string{"./foo", "co"}, w)
	expected := `fatal: missing required argument branch
usage: foo checkout [--branch] <branch>
`
	assertAppUsageOk(t, expected, w.str)
}

func assertAppUsageOk(t *testing.T, expectedOutput, actualOutput string) {
	if expectedOutput != actualOutput {
		t.Errorf("expected output: %v, found: %v", expectedOutput, actualOutput)
	}
}
