// Copyright (c) 2017. Oleg Sklyar & teris.io. All rights reserved.
// See the LICENSE file in the project root for licensing information.

package cli_test

import (
	"testing"

	"github.com/teris-io/cli"
)

func setupRunApp() cli.App {
	co := cli.NewCommand("checkout", "Check out a branch or revision").
		WithShortcut("co").
		WithArg(cli.NewArg("branch", "branch to checkout")).
		WithArg(cli.NewArg("fallback", "branch to fallback").AsOptional()).
		WithOption(cli.NewOption("branch", "Create branch if missing").WithChar('b').WithType(cli.TypeBool)).
		WithAction(func(args []string, options map[string]string) int {
			if _, ok := options["branch"]; !ok {
				return -1
			}
			if len(args) < 1 || args[0] != "5.5.5" {
				return -2
			}
			if _, ok := options["verbose"]; ok {
				return 26
			}
			if len(args) == 2 && args[1] == "5.1.1" {
				return 27
			}
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
		WithCommand(rmt).
		WithOption(cli.NewOption("verbose", "Verbose execution").WithChar('v').WithType(cli.TypeBool)).
		WithAction(func(args []string, options map[string]string) int {
			return 13
		})
}

func TestApp_Run_TopLevel_ok(t *testing.T) {
	a := setupRunApp()
	w := &stringwriter{}
	code := a.Run([]string{"./foo"}, w)
	assertAppRunOk(t, 13, code)
}

func TestApp_Run_NestedCommand_ok(t *testing.T) {
	a := setupRunApp()
	w := &stringwriter{}
	code := a.Run([]string{"./foo", "co", "-b", "5.5.5"}, w)
	assertAppRunOk(t, 25, code)
}

func TestApp_Run_NestedCommandWithOptionsFromRoot_ok(t *testing.T) {
	a := setupRunApp()
	w := &stringwriter{}
	code := a.Run([]string{"./foo", "co", "-bv", "5.5.5"}, w)
	assertAppRunOk(t, 26, code)
}

func TestApp_Run_NestedCommandWithOptionalArg_ok(t *testing.T) {
	a := setupRunApp()
	w := &stringwriter{}
	code := a.Run([]string{"./foo", "co", "-b", "5.5.5", "5.1.1"}, w)
	assertAppRunOk(t, 27, code)
}

func assertAppRunOk(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("expected exit code: %s, found: %v", expectedCode, actualCode)
	}
}
