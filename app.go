package cli

import (
	"github.com/silvertern/cli/command"
	"github.com/silvertern/cli/option"
)

type App interface {
	Description() string
	Args() []command.Arg
	Options() []option.Option
	Commands() []command.Command
	Action() command.Action

	WithArg(arg command.Arg) App
	WithOption(opt option.Option) App
	WithCommand(cmd command.Command) App
	WithAction(action command.Action) App

	Parse(appargs []string) (invocation []string, args []string, opts map[string]string, err error)
}

func New(descr string) App {
	return &app{descr: descr}
}

type app struct {
	descr  string
	args   []command.Arg
	opts   []option.Option
	cmds   []command.Command
	action command.Action
}

func (a *app) Description() string {
	return a.descr
}

func (a *app) Args() []command.Arg {
	return a.args
}

func (a *app) Options() []option.Option {
	return a.opts
}

func (a *app) Commands() []command.Command {
	return a.cmds
}

func (a *app) Action() command.Action {
	return a.action
}

func (a *app) WithArg(arg command.Arg) App {
	a.args = append(a.args, arg)
	return a
}

func (a *app) WithOption(opt option.Option) App {
	a.opts = append(a.opts, opt)
	return a
}

func (a *app) WithCommand(cmd command.Command) App {
	a.cmds = append(a.cmds, cmd)
	return a
}
func (a *app) WithAction(action command.Action) App {
	a.action = action
	return a
}

func (a *app) Parse(appargs []string) (invocation []string, args []string, opts map[string]string, err error) {
	return Parse(a, appargs)
}
