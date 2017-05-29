package cli_test

import (
	"github.com/silvertern/cli"
	"testing"
)

func setup_run_app() cli.App {
	co := cli.NewCommand("checkout", "Check out a branch or revision").
		WithShortcut("co").
		WithArg(cli.NewArg("branch", "branch to checkout")).
		WithOption(cli.NewOption("branch", "Create branch if missing").WithChar('b').WithType(cli.TypeBool)).
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

func TestApp_Run_NestedCommand_ok(t *testing.T) {
	a := setup_run_app()
	w := &stringwriter{}
	code := a.Run([]string{"./foo", "co", "-b", "5.5.5"}, w)
	assertAppRunOk(t, 25, code)
}

func assertAppRunOk(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("expected exit code: %s, found: %v", expectedCode, actualCode)
	}
}
