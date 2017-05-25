package command

import "github.com/silvertern/cli/option"

type Arg struct {
	Key      string
	Type     option.Type
	Optional bool
}
