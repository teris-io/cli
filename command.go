package cli

// Command defines a named sub-command in a command-tree of an application. A complete path to the terminal
// command e.g. `git remote add` must be defined ahead of any options or positional arguments. These are parsed
// first.
type Command interface {
	// Key returns the command name.
	Key() string
	// Shortcut returns the command shortcut if not empty.
	Shortcut() string
	// Description returns the command description to be output in the usage.
	Description() string
	// Args returns required and optional positional arguments for this command.
	Args() []Arg
	// Options permitted for this command and its sub-commands.
	Options() []Option
	// Commands returns the set of sub-commands for this command.
	Commands() []Command
	// Action returns the command action when no further sub-command is specified.
	Action() Action

	// WithShortcut adds a (shorter) command alias, e.g. `co` for `checkout`.
	WithShortcut(shortcut string) Command
	// WithArg adds a positional argument to the command. Specifying last application/command
	// argument as optional permits unlimited number of further positional arguments (at least one
	// optional argument needs to be specified in the definition for this case).
	WithArg(arg Arg) Command
	// WithOption adds a permitted option to the command and all sub-commands.
	WithOption(opt Option) Command
	// WithCommand adds a next-level sub-command to the command.
	WithCommand(cmd Command) Command
	// WithAction sets the action function for this command.
	WithAction(action Action) Command
}

// NewCommand creates a new command to be added to an application or to another command.
func NewCommand(key, descr string) Command {
	return &command{key: key, descr: descr}
}

type command struct {
	key      string
	descr    string
	shortcut string
	args     []Arg
	opts     []Option
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

func (c *command) Options() []Option {
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

func (c *command) WithOption(opt Option) Command {
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
