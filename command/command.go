package command

import (
	"github.com/silvertern/cli/option"
)

type Command interface {
	Key() string
	Description() string
	Shortcut() string
	Args() []Arg
	Options() []option.Option
	Commands() []Command
	Action() Action

	WithShortcut(shortcut string) Command
	WithArg(arg Arg) Command
	WithOption(opt option.Option) Command
	WithCommand(cmd Command) Command
	WithAction(action Action) Command
}

func New(key, descr string) Command {
	return &command{key: key, descr: descr}
}

type command struct {
	key      string
	descr    string
	shortcut string
	args     []Arg
	opts     []option.Option
	cmds     []Command
	action   Action
}

func (c *command) Key() string {
	return c.key
}

func (c *command) Description() string {
	return c.descr
}

func (c *command) Shortcut() string {
	return c.shortcut
}

func (c *command) Args() []Arg {
	return c.args
}

func (c *command) Options() []option.Option {
	return c.opts
}

func (c *command) Commands() []Command {
	return c.cmds
}

func (c *command) Action() Action {
	return c.action
}

func (c *command) WithShortcut(shortcut string) Command {
	c.shortcut = shortcut
	return c
}

func (c *command) WithArg(arg Arg) Command {
	c.args = append(c.args, arg)
	return c
}

func (c *command) WithOption(opt option.Option) Command {
	c.opts = append(c.opts, opt)
	return c
}

func (c *command) WithCommand(cmd Command) Command {
	c.cmds = append(c.cmds, cmd)
	return c
}

func (c *command) WithAction(action Action) Command {
	c.action = action
	return c
}
