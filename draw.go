package tetris

import (
	"bytes"
	"io"

	"github.com/wzshiming/ctc"
	"github.com/wzshiming/cursor"
)

var (
	wallStr = ctc.Negative.String() + "  " + ctc.Reset.String()
)

type Draw struct {
	out io.Writer
}

func NewDraw(out io.Writer) *Draw {
	return &Draw{out}
}

func (b *Draw) Clear() {
	io.WriteString(b.out, cursor.RawClear())
}

func (b *Draw) Block(block Block, point string, cw int, x, y int) {
	for i := 0; i != 4; i++ {
		for j := 0; j != 4; j++ {
			if block.On(i, j) == 1 {
				b.Dot(point, cw, x+i, y+j)
			}
		}
	}
}

func (b *Draw) Dot(point string, cw int, x, y int) {
	if x < 0 || y < 0 {
		return
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString(cursor.RawMoveTo(uint64(x*cw), uint64(y)))
	buf.WriteString(point)
	buf.WriteString("\n")
	b.out.Write(buf.Bytes())
}

func (b *Draw) Box(point string, cw int, x, y, w, h int) {
	buf := bytes.NewBuffer(nil)
	x -= 1
	y -= 1
	buf.WriteString(cursor.RawMoveTo(uint64(x*cw), uint64(y)))
	for i := 0; i <= w+1; i++ {
		buf.WriteString(point)
	}
	for i := 0; i < h; i++ {
		buf.WriteString(cursor.RawMoveTo(uint64(x*cw), uint64(y+i+1)))
		buf.WriteString(point)
		buf.WriteString(cursor.RawMoveTo(uint64((x+w+1)*cw), uint64(y+i+1)))
		buf.WriteString(point)
	}
	buf.WriteString(cursor.RawMoveTo(uint64(x*cw), uint64(y+h+1)))
	for i := 0; i <= w+1; i++ {
		buf.WriteString(point)
	}
	buf.WriteString("\n")
	b.out.Write(buf.Bytes())
}
