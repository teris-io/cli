package option

type Type int

const (
	TypeString Type = iota
	TypeBool
	TypeInt
	TypeNumber
)

type Option interface {
	Key() string
	CharKey() rune
	Description() string
	Type() Type
	WithChar(char rune) Option
	WithType(ft Type) Option
}

type option struct {
	key   string
	char  rune
	descr string
	tp    Type
}

func New(key, descr string) Option {
	char := rune(0)
	if len(key) == 1 {
		char = rune(key[0])
	}
	return option{key: key, char: char, descr: descr, tp: TypeString}
}

func (f option) Key() string {
	return f.key
}

func (f option) CharKey() rune {
	return f.char
}

func (f option) Type() Type {
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

func (f option) WithType(tp Type) Option {
	f.tp = tp
	return f
}
