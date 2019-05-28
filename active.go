package tetris

import (
	"github.com/wzshiming/ctc"
)

var (
	wallStr   = ctc.Negative.String() + "  " + ctc.Reset.String()
	bottomStr = ctc.Negative.String() + "                        " + ctc.Reset.String()
)

var Actives = []string{
	ctc.BackgroundBright.String() + "  " + ctc.Reset.String(),
	ctc.BackgroundRed.String() + "  " + ctc.Reset.String(),
	ctc.BackgroundGreen.String() + "  " + ctc.Reset.String(),
	ctc.BackgroundYellow.String() + "  " + ctc.Reset.String(),
	ctc.BackgroundBlue.String() + "  " + ctc.Reset.String(),
	ctc.BackgroundMagenta.String() + "  " + ctc.Reset.String(),
	ctc.BackgroundCyan.String() + "  " + ctc.Reset.String(),
}
