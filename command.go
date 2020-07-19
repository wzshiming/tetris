package tetris

type Command uint

const (
	None Command = iota
	Pause
	RightRotate
	LeftRotate
	RightMove
	LeftMove
	DownMove
	Drop
)
