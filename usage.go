// Copyright (c) 2017. Oleg Sklyar & teris.io. All rights reserved.
// See the LICENSE file in the project root for licensing information.

package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type usageline struct {
	section string
	key     string
	value   string
}

// Usage prints out the complete usage string.
func Usage(a App, invocation []string, w io.Writer) error {
	if len(invocation) < 1 {
		return errors.New("invalid invocation path []")
	}

	descr := a.Description()
	cmds := a.Commands()
	args := a.Args()
	opts := a.Options()

	if len(invocation) > 1 {
		for _, key := range invocation[1:] {
			matched := false
			for _, cmd := range cmds {
				if cmd.Key() == key {
					descr = cmd.Description()
					cmds = cmd.Commands()
					args = cmd.Args()
					opts = append(opts, cmd.Options()...)
					matched = true
					break
				}
			}
			// should never happen if invocation originates from the parser
			if !matched {
				// ignore errors here as no alternative writer is available
				err := fmt.Errorf("fatal: invalid invocation path %v\n", invocation)
				fmt.Fprint(w, err.Error())
				return err
			}
		}
	}

	indent := "    "
	thiscmd := strings.Join(invocation, " ")
	fmt.Fprintf(w, "%s%s%s\n\n", thiscmd, optstring(opts), argstring(args))
	fmt.Fprintln(w, "Description:")
	fmt.Fprintf(w, "%s%s\n", indent, descr)

	var lines []usageline
	maxkey := 0
	if len(args) > 0 {
		for _, arg := range args {
			value := arg.Description()
			if arg.Optional() {
				value += ", optional"
			}
			line := usageline{
				section: "Arguments",
				key:     arg.Key(),
				value:   value,
			}
			lines = append(lines, line)
			if len(line.key) > maxkey {
				maxkey = len(line.key)
			}
		}
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			charstr := "    "
			if opt.CharKey() != rune(0) {
				charstr = "-" + string(opt.CharKey()) + ", "
			}

			line := usageline{
				section: "Options",
				key:     charstr + "--" + opt.Key(),
				value:   opt.Description(),
			}
			lines = append(lines, line)
			if len(line.key) > maxkey {
				maxkey = len(line.key)
			}
		}
	}

	if len(cmds) > 0 {
		for _, cmd := range cmds {
			shortstr := ""
			if cmd.Shortcut() != "" {
				shortstr = ", shortcut: " + cmd.Shortcut()
			}

			line := usageline{
				section: "Sub-commands",
				key:     thiscmd + " " + cmd.Key(),
				value:   cmd.Description() + shortstr,
			}
			lines = append(lines, line)
			if len(line.key) > maxkey {
				maxkey = len(line.key)
			}
		}
	}

	lastsection := ""
	for _, line := range lines {
		if line.section != lastsection {
			fmt.Fprintf(w, "\n%s:\n", line.section)
		}
		lastsection = line.section
		spacer := 3 + maxkey - len(line.key)
		fmt.Fprintf(w, "%s%s%s%s\n", indent, line.key, strings.Repeat(" ", spacer), line.value)
	}
	return nil
}

func optstring(opts []Option) string {
	res := ""
	for _, opt := range opts {
		res += " [--" + opt.Key()
		switch opt.Type() {
		case TypeString:
			res += "=string"
		case TypeInt:
			res += "=int"
		case TypeNumber:
			res += "=number"
		}
		res += "]"
	}
	return res
}

func argstring(args []Arg) string {
	res := ""
	for _, arg := range args {
		if arg.Optional() {
			res += " [" + arg.Key() + "]"
		} else {
			res += " <" + arg.Key() + ">"
		}
	}
	return res
}

func shortUsage(a App, invocation []string) string {
	if len(invocation) < 1 {
		return "invalid invocation path []"
	}

	cmds := a.Commands()
	args := a.Args()
	opts := a.Options()

	if len(invocation) > 1 {
		for _, key := range invocation[1:] {
			matched := false
			for _, cmd := range cmds {
				if cmd.Key() == key {
					cmds = cmd.Commands()
					args = cmd.Args()
					opts = append(opts, cmd.Options()...)
					matched = true
					break
				}
			}
			// should never happen if invocation originates from the parser
			if !matched {
				return fmt.Sprintf("fatal: invalid invocation path %v", invocation)
			}
		}
	}

	return fmt.Sprintf("%s%s%s", strings.Join(invocation, " "), optstring(opts), argstring(args))
}
