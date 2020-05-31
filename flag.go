package tetris

type flag uint8

const (
	_ flag = (8 - iota) << 1
	overFlag
	pauseFlag
)

func (f *flag) On(a flag) {
	*f |= a
}

func (f *flag) Off(a flag) {
	*f &^= a
}

func (f *flag) Has(a flag) bool {
	return (*f)&a == a
}
