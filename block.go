package tetris

import (
	"github.com/wzshiming/ctc"
)

type Block uint16

func (b Block) On(x, y int) int {
	if x > 3 || y > 3 || x < 0 || y < 0 {
		return 0
	}
	off := 4*(3-y) + (3 - x)
	if b&(1<<uint(off)) != 0 {
		return 1
	}
	return 0
}

type Blocks struct {
	Blocks  []Block
	Color   string
	Predict string
}

var (
	blockRune   = "  "
	predictRune = "<>"
)
var BlocksPool = [...]Blocks{
	{
		Blocks: []Block{
			// 0 0 0 0
			// 1 1 1 1
			// 0 0 0 0
			// 0 0 0 0
			0b_0000_1111_0000_0000,

			// 0 0 1 0
			// 0 0 1 0
			// 0 0 1 0
			// 0 0 1 0
			0b_0010_0010_0010_0010,

			// 0 0 0 0
			// 0 0 0 0
			// 1 1 1 1
			// 0 0 0 0
			0b_0000_0000_1111_0000,

			// 0 1 0 0
			// 0 1 0 0
			// 0 1 0 0
			// 0 1 0 0
			0b_0100_0100_0100_0100,
		},
		Color:   ctc.BackgroundBright.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundBright.String() + predictRune + ctc.Reset.String(),
	},

	{
		Blocks: []Block{
			// 0 0 0 0
			// 1 1 0 0
			// 0 1 1 0
			// 0 0 0 0
			0b_0000_1100_0110_0000,

			// 0 0 0 0
			// 0 0 1 0
			// 0 1 1 0
			// 0 1 0 0
			0b_0000_0010_0110_0100,

			// 0 0 0 0
			// 0 0 0 0
			// 1 1 0 0
			// 0 1 1 0
			0b_0000_0000_1100_0110,

			// 0 0 0 0
			// 0 1 0 0
			// 1 1 0 0
			// 1 0 0 0
			0b_0000_0100_1100_1000,
		},
		Color:   ctc.BackgroundRed.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundRed.String() + predictRune + ctc.Reset.String(),
	},

	{
		Blocks: []Block{
			// 0 0 0 0
			// 0 1 1 0
			// 1 1 0 0
			// 0 0 0 0
			0b_0000_0110_1100_0000,

			// 0 0 0 0
			// 0 1 0 0
			// 0 1 1 0
			// 0 0 1 0
			0b_0000_0100_0110_0010,

			// 0 0 0 0
			// 0 0 0 0
			// 0 1 1 0
			// 1 1 0 0
			0b_0000_0000_0110_1100,

			// 0 0 0 0
			// 1 0 0 0
			// 1 1 0 0
			// 0 1 0 0
			0b_0000_1000_1100_0100,
		},
		Color:   ctc.BackgroundGreen.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundGreen.String() + predictRune + ctc.Reset.String(),
	},

	{
		Blocks: []Block{
			// 0 0 0 0
			// 1 0 0 0
			// 1 1 1 0
			// 0 0 0 0
			0b_0000_1000_1110_0000,

			// 0 0 0 0
			// 0 1 1 0
			// 0 1 0 0
			// 0 1 0 0
			0b_0000_0110_0100_0100,

			// 0 0 0 0
			// 0 0 0 0
			// 1 1 1 0
			// 0 0 1 0
			0b_0000_0000_1110_0010,

			// 0 0 0 0
			// 0 1 0 0
			// 0 1 0 0
			// 1 1 0 0
			0b_0000_0100_0100_1100,
		},
		Color:   ctc.BackgroundBlue.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundBlue.String() + predictRune + ctc.Reset.String(),
	},

	{
		Blocks: []Block{
			// 0 0 0 0
			// 0 0 1 0
			// 1 1 1 0
			// 0 0 0 0
			0b_0000_0010_1110_0000,

			// 0 0 0 0
			// 0 1 0 0
			// 0 1 0 0
			// 0 1 1 0
			0b_0000_0100_0100_0110,

			// 0 0 0 0
			// 0 0 0 0
			// 1 1 1 0
			// 1 0 0 0
			0b_0000_0000_1110_1000,

			// 0 0 0 0
			// 1 1 0 0
			// 0 1 0 0
			// 0 1 0 0
			0b_0000_1100_0100_0100,
		},
		Color:   ctc.BackgroundCyan.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundCyan.String() + predictRune + ctc.Reset.String(),
	},

	{
		Blocks: []Block{
			// 0 0 0 0
			// 0 1 0 0
			// 1 1 1 0
			// 0 0 0 0
			0b_0000_0100_1110_0000,

			// 0 0 0 0
			// 0 1 0 0
			// 0 1 1 0
			// 0 1 0 0
			0b_0000_0100_0110_0100,

			// 0 0 0 0
			// 0 0 0 0
			// 1 1 1 0
			// 0 1 0 0
			0b_0000_0000_1110_0100,

			// 0 0 0 0
			// 0 1 0 0
			// 1 1 0 0
			// 0 1 0 0
			0b_0000_0100_1100_0100,
		},
		Color:   ctc.BackgroundMagenta.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundMagenta.String() + predictRune + ctc.Reset.String(),
	},

	{
		Blocks: []Block{
			// 0 0 0 0
			// 0 1 1 0
			// 0 1 1 0
			// 0 0 0 0
			0b_0000_0110_0110_0000,
		},
		Color:   ctc.BackgroundYellow.String() + blockRune + ctc.Reset.String(),
		Predict: ctc.ForegroundYellow.String() + predictRune + ctc.Reset.String(),
	},
}
