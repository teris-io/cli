package cli

import (
	"fmt"
	"github.com/silvertern/cli/command"
	"github.com/silvertern/cli/option"
	"path"
	"strconv"
	"strings"
)

type App struct {
	Description string
	Args        []command.Arg
	Options     []option.Option
	Commands    []command.Command
	Action      command.Action
}

func (a *App) Parse(appargs []string) (invocation []string, args []string, opts map[string]string, err error) {
	opts = make(map[string]string)
	args = []string{}

	_, appname := path.Split(appargs[0])
	invocation = []string{appname}

	appargs = appargs[1:]

	permittedOpts := a.Options
	expectedArgs := a.Args
	availableCmds := a.Commands
	for _, apparg := range appargs {
		matched := false
		for _, cmd := range availableCmds {
			if cmd.Key() == apparg || cmd.Shortcut() == apparg {
				invocation = append(invocation, cmd.Key())
				permittedOpts = append(permittedOpts, cmd.Options()...)
				expectedArgs = cmd.Args()
				appargs = appargs[1:]
				availableCmds = cmd.Commands()
				matched = true
				break
			}
		}
		if !matched {
			break
		}
	}

	var expectingValueFor string
	for _, apparg := range appargs {
		if expectingValueFor != "" {
			opts[expectingValueFor] = apparg
			expectingValueFor = ""
		} else if strings.HasPrefix(apparg, "--") {
			apparg = apparg[2:]
			if apparg == "help" {
				return invocation, nil, map[string]string{"help": "true"}, nil
			}
			parts := strings.Split(apparg, "=")
			matched := false
			for _, permittedOpt := range permittedOpts {
				if permittedOpt.Key() == parts[0] {
					if permittedOpt.Type() == option.TypeBool {
						opts[permittedOpt.Key()] = "true"
					} else if len(parts) != 2 {
						return invocation, args, opts, fmt.Errorf("missing value for option --%s", permittedOpt.Key())
					} else {
						opts[permittedOpt.Key()] = parts[1]
					}
					matched = true
					break
				}
			}
			if ! matched {
				return invocation, args, opts, fmt.Errorf("unknown option --%s", parts[0])
			}
		} else if strings.HasPrefix(apparg, "-") {
			apparg = apparg[1:]

			for i, char := range apparg {
				if string(char) == "h" {
					return invocation, nil, map[string]string{"help": "true"}, nil
				}
				matched := false
				for _, permittedOpt := range permittedOpts {
					if permittedOpt.CharKey() == char {
						if permittedOpt.Type() == option.TypeBool {
							opts[permittedOpt.Key()] = "true"
						} else if i == len(apparg)-1 {
							expectingValueFor = permittedOpt.Key()
						} else {
							return invocation, args, opts, fmt.Errorf("non-boolean flag -%v in non-terminal position", string(char))
						}
						matched = true
						break
					}
				}
				if !matched {
					return invocation, args, opts, fmt.Errorf("unknown flag -%v", string(char))
				}
			}
		} else {
			args = append(args, apparg)
		}
	}
	if expectingValueFor != "" {
		return invocation, args, opts, fmt.Errorf("dangling option --%s", expectingValueFor)
	}

	lastArgOptional := false
	if len(expectedArgs) > 0 && expectedArgs[len(expectedArgs) - 1].Optional {
		lastArgOptional = true
	}
	if !lastArgOptional {
		if len(expectedArgs) > len(args) {
			return invocation, args, opts, fmt.Errorf("missing required argument %v", expectedArgs[len(args)].Key)
		}	else if len(expectedArgs) < len(args) {
			return invocation, args, opts, fmt.Errorf("unknown arguments %v", args[len(expectedArgs):])
		}
	}
	for i, expectedArg := range expectedArgs {
		if len(args) < i+1 {
			if expectedArg.Optional {
				break
			}
			return invocation, args, opts, fmt.Errorf("missing required argument %s", expectedArg.Key)
		}
		arg := args[i]
		switch expectedArg.Type {
		case option.TypeBool:
			if _, err := strconv.ParseBool(arg); err != nil {
				return invocation, args, opts, fmt.Errorf("argument %s must be a boolean value, found %v", expectedArg.Key, arg)
			}
		case option.TypeInt:
			if _, err := strconv.ParseInt(arg, 10, 64); err != nil {
				return invocation, args, opts, fmt.Errorf("argument %s must be an integer value, found %v", expectedArg.Key, arg)
			}
		case option.TypeNumber:
			if _, err := strconv.ParseFloat(arg, 64); err != nil {
				return invocation, args, opts, fmt.Errorf("argument %s must be a number, found %v", expectedArg.Key, arg)
			}
		default:
		}
	}

	for key, value := range opts {
		for _, permittedOpt := range permittedOpts {
			if permittedOpt.Key() == key {
				switch permittedOpt.Type() {
				case option.TypeInt:
					if _, err := strconv.ParseInt(value, 10, 64); err != nil {
						return invocation, args, opts, fmt.Errorf("option --%s must be given an integer value, found %v", permittedOpt.Key(), value)
					}
				case option.TypeNumber:
					if _, err := strconv.ParseFloat(value, 64); err != nil {
						return invocation, args, opts, fmt.Errorf("option --%s must must be given a number, found %v", permittedOpt.Key(), value)
					}
				default:
				}
				break
			}
		}
	}

	return invocation, args, opts, nil
}
