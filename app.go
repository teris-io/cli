// Package cli provides a simple, fast and complete API for building command line applications in Go.
// In contrast to other libraries additional emphasis is put on the definition and validation of
// positional arguments and consistent usage outputs combining options from all command levels into
// one block.
package cli

// Action defines a function type to be executed for an application or a
// command. It takes a slice of validated positional arguments and a map
// of validated options (with all value types encoded as strings) and
// returns a Unix exit code (success: 0).
type Action func(args []string, options map[string]string) int

// App defines a CLI application parameterizable with sub-commands, arguments and options.
type App interface {
	// Description returns the application description to be output in the usage.
	Description() string
	// Args returns required and optional positional arguments for the top-level application.
	Args() []Arg
	// Options permitted for the top-level application and all sub-commands.
	Options() []Option
	// Commands returns the set of first-level sub-commands for the application.
	Commands() []Command
	// Action returns the application action when no sub-command is specified.
	Action() Action

	// WithArg adds a positional argument to the application. Specifying last application/command
	// argument as optional permits unlimited number of further positional arguments (at least one
	// optional argument needs to be specified in the definition for this case).
	WithArg(arg Arg) App
	// WithOption adds a permitted option to the application and all sub-commands.
	WithOption(opt Option) App
	// WithCommand adds a first-level sub-command to the application.
	WithCommand(cmd Command) App
	// WithAction sets the action function to execute after successful parsing of commands, arguments
	// and options to the top-level application.
	WithAction(action Action) App

	// Parse parses the original application arguments into the command invocation path (application ->
	// first level command -> second level command etc.), a list of validated positional arguments matching
	// the command being invoked (the last one in the invocation path) and a map of validated options
	// matching one of the invocation path elements, from the top application down to the command being invoked.
	// An error is returned if a command is not found or arguments or options are invalid. In case of an error,
	// the invocation path is normally also computed and returned (the content of arguments and options is not
	// guaranteed).
	Parse(appargs []string) (invocation []string, args []string, opts map[string]string, err error)
}

// New creates a new CLI App.
func New(descr string) App {
	return &app{descr: descr}
}

// ValueType defines the type of permitted argument and option values.
type ValueType int

// ValueType constants for string, boolean, int and number options and arguments.
const (
	TypeString ValueType = iota
	TypeBool
	TypeInt
	TypeNumber
)

// Arg defines a positional argument. Arguments are validated for their
// count and their type. If the last defined argument is optional, then
// an unlimited number of arguments can be passed into the call, otherwise
// an exact count of positional arguments is expected. Obviously, optional
// arguments can be omitted. No validation is done for the invalid case of
// specifying an optional positional argument before a required one.
type Arg interface {
	// Key defines how the argument will be shown in the usage string.
	Key() string
	// Type defines argument type. Default is string, which is not validated,
	// other types are validated by simple string parsing into boolean, int and float.
	Type() ValueType
	// Optional specifies that an argument may be omitted. No non-optional arguments
	// should follow an optional one (no validation for this scenario as this is
	// the definition time exception, rather than incorrect input at runtime).
	Optional() bool

	// WithType sets the argument type.
	WithType(at ValueType) Arg
	// AsOptional sets the argument as optional.
	AsOptional() Arg
}

// NewArg creates a new positional argument.
func NewArg(key string) Arg {
	return arg{key: key}
}

type app struct {
	descr  string
	args   []Arg
	opts   []Option
	cmds   []Command
	action Action
}

func (a *app) Description() string {
	return a.descr
}

func (a *app) Args() []Arg {
	return a.args
}

func (a *app) Options() []Option {
	return a.opts
}

func (a *app) Commands() []Command {
	return a.cmds
}

func (a *app) Action() Action {
	return a.action
}

func (a *app) WithArg(arg Arg) App {
	a.args = append(a.args, arg)
	return a
}

func (a *app) WithOption(opt Option) App {
	a.opts = append(a.opts, opt)
	return a
}

func (a *app) WithCommand(cmd Command) App {
	a.cmds = append(a.cmds, cmd)
	return a
}
func (a *app) WithAction(action Action) App {
	a.action = action
	return a
}

func (a *app) Parse(appargs []string) (invocation []string, args []string, opts map[string]string, err error) {
	return Parse(a, appargs)
}

type arg struct {
	key      string
	at       ValueType
	optional bool
}

func (a arg) Key() string {
	return a.key
}

func (a arg) Type() ValueType {
	return a.at
}

func (a arg) Optional() bool {
	return a.optional
}

func (a arg) WithType(at ValueType) Arg {
	a.at = at
	return a
}

func (a arg) AsOptional() Arg {
	a.optional = true
	return a
}
