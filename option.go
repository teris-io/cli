// Copyright (c) 2017. Oleg Sklyar & teris.io. All rights reserved.
// See the LICENSE file in the project root for licensing information.

package cli

// Option defines an application or command option.
//
// Boolean options do not need a value, as their presence implicitly means `true`. All other option
// types need a value. When options are specified using their char keys their need to be prepended
// by a single dash (-) and can be joined together (-zxv). Hhowever, only boolean options requiring
// no further value may occupy a non-terminal position of a join. A value must follow a char option
// as the next argument, same applies for a non-boolean option at the terminal position of an option
// join, e.g. `-fc 1` means `-f -c 1`, where 1 is the argument for `-c`.
//
// Non-char, complete, options must be prepended with a double dash (--) and their value must be
// provided in the same argument after the equal sign, e.g. --count=1. Similarly here, boolean options
// require no value. Empty values are supported for complete (non-char) string options only by providing
// no value after the equal sign, e.g. `--default=`. Equal signs can be used within the option value, e.g.
// `--default=a=b=c` specifies the `a=b=c` string as a value for `--default`.
//
// Every option must have a `complete` name as these are used as keys to pass options to the action. In case
// only a char option is desired, a complete key with the same single char should be defined.
//
// Options can be used at any position after the command, arbitrarily intermixed with positional arguments.
// In contrast to positional arguments the order of options is not preserved.
type Option interface {
	// Key returns the complete key of an option (used with the -- notation), required.
	Key() string
	// CharKey returns a single-character key of an option (used with the - notation), optional.
	CharKey() rune
	// Description returns the option description for the usage string.
	Description() string
	// Type returns the option type (string by default) to be used to decide if a value is required and for
	// value validation.
	Type() ValueType

	// WithChar sets the char key for the option.
	WithChar(char rune) Option
	// WithType sets the option value type.
	WithType(ft ValueType) Option
}

// NewOption creates a new option with a given key and description.
func NewOption(key, descr string) Option {
	char := rune(0)
	if len(key) == 1 {
		char = rune(key[0])
	}
	return option{key: key, char: char, descr: descr, tp: TypeString}
}

type option struct {
	key   string
	char  rune
	descr string
	tp    ValueType
}

func (f option) Key() string {
	return f.key
}

func (f option) CharKey() rune {
	return f.char
}

func (f option) Type() ValueType {
	return f.tp
}

func (f option) Description() string {
	return f.descr
}

func (f option) String() string {
	return f.descr
}

func (f option) WithChar(char rune) Option {
	f.char = char
	return f
}

func (f option) WithType(tp ValueType) Option {
	f.tp = tp
	return f
}
