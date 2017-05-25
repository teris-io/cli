package command

type Action func(args []string, options map[string]string) int
