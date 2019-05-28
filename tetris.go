package tetris

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/wzshiming/cursor"
	"github.com/wzshiming/getch"
)

type Tetris struct {
	box           [20][10]string
	x, y          int8
	waitingRotate uint8
	waiting       []Block
	waitingColor  string
	currentRotate uint8
	current       []Block
	currentColor  string
	emptyColor    string
	rand          *rand.Rand
	rank          uint64
	over          uint8
}

func NewTetris() *Tetris {
	t := &Tetris{
		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
	t.init()
	t.next()
	return t
}

func (t *Tetris) Show() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(cursor.RawClear())
	buf.WriteString(cursor.RawMoveUp(uint64(len(t.box) + 1)))

	for j, row := range t.box {
		buf.WriteString(wallStr)
		for i := range row {
			buf.WriteString(t.Get(int8(i), int8(j)))
		}
		buf.WriteString(wallStr)
		buf.WriteString("\n")
	}
	buf.WriteString(bottomStr)
	buf.WriteString("\n")
	return buf.String()
}

func (t *Tetris) init() {
	t.emptyColor = "  "
	t.waiting = BlocksPool[t.rand.Int()%len(BlocksPool)]
	t.waitingColor = Actives[t.rand.Int()%len(Actives)]
	t.waitingRotate = uint8(t.rand.Int() % len(t.waiting))
}

func (t *Tetris) next() {
	t.y = -4
	t.x = 3

	t.current = t.waiting
	t.currentColor = t.waitingColor
	t.currentRotate = t.waitingRotate
	t.waiting = BlocksPool[t.rand.Int()%len(BlocksPool)]
	t.waitingColor = Actives[t.rand.Int()%len(Actives)]
	t.waitingRotate = uint8(t.rand.Int() % len(t.waiting))
}

func (t *Tetris) Get(x, y int8) string {
	if x > 10 || y > 20 || x < 0 || y < 0 {
		return t.emptyColor
	}
	nx := x - t.x
	ny := y - t.y

	if nx >= 0 && nx < 4 && ny >= 0 && ny < 4 {
		block := t.current[t.currentRotate]
		if block.On(nx, ny) == 1 {
			return t.currentColor
		}
	}
	b := t.box[y][x]
	if b == "" {
		return t.emptyColor
	}
	return b
}

func (t *Tetris) Set(x, y int8, currentColor string) {
	if x > 10 || y > 20 || x < 0 || y < 0 {
		t.over = 1
		return
	}
	t.box[y][x] = currentColor
}

func (t *Tetris) On(x, y int8) int8 {
	if x >= 10 || y >= 20 || x < 0 {
		return 1
	}
	if y < 0 {
		return 0
	}
	if t.box[y][x] != "" {
		return 1
	}
	return 0
}

func (t *Tetris) touch(block Block, x, y int8) int8 {
	for i := int8(0); i != 4; i++ {
		for j := int8(0); j != 4; j++ {
			if block.On(i, j) == 1 && t.On(x+i, y+j) == 1 {
				return 1
			}
		}
	}
	return 0
}

func (t *Tetris) eliminate(y int8) {
	if y >= 20 || y < 0 {
		return
	}
	eli := true
	for _, d := range t.box[y] {
		if d == "" {
			eli = false
			break
		}
	}
	if eli {
		t.rank++
		copy(t.box[1:y+1], t.box[:y])
		t.box[0] = [10]string{}
	}
}

func (t *Tetris) merge(block Block, x, y int8) {
	for j := int8(0); j != 4; j++ {
		for i := int8(0); i != 4; i++ {
			if block.On(i, j) == 1 {
				t.Set(x+i, y+j, t.currentColor)
			}
		}
		t.eliminate(y + j)
	}
}

func (t *Tetris) Rotate() {
	currentRotate := t.currentRotate + 1
	if currentRotate >= uint8(len(t.current)) {
		currentRotate = 0
	}
	if t.touch(t.current[currentRotate], t.x, t.y) == 0 {
		t.currentRotate = currentRotate
	}
}

func (t *Tetris) Left() {
	if t.touch(t.current[t.currentRotate], t.x-1, t.y) == 0 {
		t.x--
	}
}

func (t *Tetris) Right() {
	if t.touch(t.current[t.currentRotate], t.x+1, t.y) == 0 {
		t.x++
	}
}

func (t *Tetris) Drop() {
	i := 1
	for i != 0 {
		i = t.down()
	}
}

func (t *Tetris) Down() {
	t.down()
}

func (t *Tetris) down() int {
	if t.touch(t.current[t.currentRotate], t.x, t.y+1) == 0 {
		t.y++
		return 1
	}
	t.merge(t.current[t.currentRotate], t.x, t.y)
	t.next()
	return 0
}

func (t *Tetris) Run() {

	type Command uint

	const (
		None Command = iota
		Rotate
		Right
		Left
		Down
		Drop
		Exit
	)

	tick := time.NewTicker(time.Second)
	cch := make(chan Command, 0)

	go func() {
		for {
			b, _, err := getch.Getch()
			if err != nil {
				cch <- Exit
			}
			c := None

			switch b {
			case 'w', 'W':
				c = Rotate
			case 's', 'S':
				c = Down
			case 'a', 'A':
				c = Left
			case 'd', 'D':
				c = Right
			case ' ':
				c = Drop
			case 'q', 'Q':
				c = Exit
			default:
				continue
			}
			select {
			case cch <- c:
			default:
				return
			}
		}
	}()

	go func() {
		for range tick.C {
			cch <- Down
		}
	}()

	for c := range cch {
		if t.over == 1 {
			c = Exit
		}
		switch c {
		case Rotate:
			t.Rotate()
		case Right:
			t.Right()
		case Left:
			t.Left()
		case Drop:
			t.Drop()
		case Down:
			t.Down()
		case Exit:
			close(cch)
			tick.Stop()
			return
		}
		io.WriteString(os.Stdout, t.Show())
	}
}
