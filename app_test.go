package cli_test

import (
	"fmt"
	"github.com/silvertern/cli"
	"github.com/silvertern/cli/command"
	"github.com/silvertern/cli/option"
	"sort"
	"testing"
)

func setup() *cli.App {
	co := command.New("checkout", "checkout a branch or revision").
		WithShortcut("co").
		WithArg(command.Arg{Key: "branch"}).
		WithOption(option.New("branch", "Create branch").WithChar('b').WithType(option.TypeBool)).
		WithOption(option.New("upstream", "Set upstream").WithChar('u').WithType(option.TypeBool)).
		WithOption(option.New("fallback", "Set upstream").WithChar('f')).
		WithOption(option.New("count", "Count").WithChar('c').WithType(option.TypeInt)).
		WithOption(option.New("pi", "Set upstream").WithChar('p').WithType(option.TypeNumber))

	add := command.New("add", "add a remote").
		WithArg(command.Arg{Key: "remote"}).
		WithArg(command.Arg{Key: "count", Type: option.TypeInt}).
		WithArg(command.Arg{Key: "pi", Type: option.TypeNumber}).
		WithArg(command.Arg{Key: "force", Type: option.TypeBool}).
		WithArg(command.Arg{Key: "optional", Optional: true}).
		WithOption(option.New("force", "Force").WithChar('f').WithType(option.TypeBool)).
		WithOption(option.New("quiet", "Quiet").WithChar('q').WithType(option.TypeBool)).
		WithOption(option.New("default", "Default"))

	rmt := command.New("remote", "operations with remotes").WithCommand(add)

	return &cli.App{
		Commands: []command.Command{
			co, rmt,
		},
	}
}

func TestApp_Parse_NoFlags_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[]", invocation, args, opts, err)
}

func TestApp_Parse_1xCharBoolFlag_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-b", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xCharBoolFlags_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-b", "-u", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xCharBoolFlagsAsOne_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-bu", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_MultiCharStringLast_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-buf", "master", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true fallback:master upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_MultiCharIntLast_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-buc", "1", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true count:1 upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_MultiCharNumberLast_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-bup", "3.14", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true pi:3.14 upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_1xBoolFlag_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "--branch", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xBoolFlag_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "--branch", "--upstream", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_2xBoolAnd1xStringFlag_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "--fallback=master", "--branch", "--upstream", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true fallback:master upstream:true]", invocation, args, opts, err)
}

func TestApp_Parse_RedundantFlags_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "checkout", "-b", "--branch", "dev"})
	assertAppParseOk(t, "[git checkout] [dev] map[branch:true]", invocation, args, opts, err)
}

func TestApp_Parse_NestedCommandWithFlags_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "remote", "add", "origin", "-f", "1", "3.14", "true"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[force:true]", invocation, args, opts, err)
}

func TestApp_Parse_OptionalMissing_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "remote", "add", "origin", "1", "3.14", "true"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[]", invocation, args, opts, err)
}

func TestApp_Parse_OptionalPresent_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "remote", "add", "origin", "1", "3.14", "true", "stuff"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true stuff] map[]", invocation, args, opts, err)
}

func TestApp_Parse_KeysAnywhereBetweenArgs_Ok(t *testing.T) {
	invocation, args, opts, err := setup().Parse([]string{"git", "remote", "add", "-f", "origin", "--default=foo", "1", "3.14", "true", "-q"})
	assertAppParseOk(t, "[git remote add] [origin 1 3.14 true] map[default:foo force:true quiet:true]", invocation, args, opts, err)
}

func assertAppParseOk(t *testing.T, expected string, invocation []string, args []string, opts map[string]string, err error) {
	if err == nil {
		optkeys := []string{}
		for key, _ := range opts {
			optkeys = append(optkeys, key)
		}
		sort.Strings(optkeys)
		for i, key := range optkeys {
			optkeys[i] = fmt.Sprintf("%s:%s", key, opts[key])
		}
		actual := fmt.Sprintf("%v %v map%v", invocation, args, optkeys)
		if actual != expected {
			t.Errorf("assertion error: expected '%v', found '%v'", expected, actual)
		}
	} else {
		t.Errorf("no error expected, found '%v'; data %v %v %v", err, invocation, args, opts)
	}
}
