package tetris

type Command uint

const (
	_ Command = iota
	Pause
	RightRotate
	LeftRotate
	RightMove
	LeftMove
	DownMove
	Drop
)
