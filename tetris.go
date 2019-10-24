package tetris

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/wzshiming/getch"
)

type Tetris struct {
	offX, offY    int
	start         time.Time
	draw          *Draw
	box           [20][10]string
	x, y          int
	waiting       Blocks
	currentRotate int8
	current       Blocks
	emptyColor    string
	rand          *rand.Rand
	rank          uint64
	over          uint
}

func NewTetris() *Tetris {
	t := &Tetris{
		start: time.Now(),
		draw:  NewDraw(os.Stdout),
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
	t.init()
	t.initBox()
	t.next()
	return t
}

func (t *Tetris) end() {
	t.draw.Dot("GAME OVER", 2, t.offX, t.offY+21)
}

func (t *Tetris) initBox() {
	t.draw.Clear()
	t.setRank(t.rank)
	t.setTime()
	t.offX = 2
	t.offY = 3
	helpX := 14
	helpY := 15
	t.draw.Dot("  H E L P", 2, helpX, helpY)
	t.draw.Dot("L:     Quit", 2, helpX, helpY+1)
	t.draw.Dot("Q:     Left rotate", 2, helpX, helpY+2)
	t.draw.Dot("E:     Right rotate", 2, helpX, helpY+3)
	t.draw.Dot("A:     Left move", 2, helpX, helpY+4)
	t.draw.Dot("D:     Right move", 2, helpX, helpY+5)
	t.draw.Dot("S:     Down", 2, helpX, helpY+6)
	t.draw.Dot("W:     Drop", 2, helpX, helpY+7)
	t.draw.Box(wallStr, 2, t.offX+13, t.offY, 4, 4)
	t.draw.Box(wallStr, 2, t.offX, t.offY, 10, 20)
}

func (t *Tetris) showBlock(block Block, point string, cw int, x, y int) {
	for i := 0; i != 4; i++ {
		for j := 0; j != 4; j++ {
			if (y+j >= 0) && block.On(i, j) == 1 {
				t.draw.Dot(point, cw, t.offX+x+i, t.offY+y+j)
			}
		}
	}
}

func (t *Tetris) showRow(row [10]string, cw int, x, y int) {
	for i, col := range row {
		if col == "" {
			t.draw.Dot(t.emptyColor, cw, t.offX+x+i, t.offY+y)
		} else {
			t.draw.Dot(col, cw, t.offX+x+i, t.offY+y)
		}
	}
}

func (t *Tetris) setRank(i uint64) {
	t.rank = i
	t.draw.Dot(fmt.Sprintf("Rank:  %d", i), 2, 14, 11)
}

func (t *Tetris) setTime() {
	t.draw.Dot(fmt.Sprintf("Time:  %s", time.Now().Sub(t.start)/time.Second*time.Second), 2, 14, 12)
}

func (t *Tetris) init() {
	t.emptyColor = "  "
	t.waiting = BlocksPool[t.rand.Int()%len(BlocksPool)]
}

func (t *Tetris) next() {
	t.y = -4
	t.x = 3

	t.current = t.waiting
	t.currentRotate = 0

	t.showBlock(t.waiting.Blocks[0], t.emptyColor, 2, int(t.offX+11), int(t.y+4))
	t.waiting = BlocksPool[t.rand.Int()%len(BlocksPool)]
	t.showBlock(t.waiting.Blocks[0], t.waiting.Color, 2, int(t.offX+11), int(t.y+4))
}

func (t *Tetris) Get(x, y int) string {
	if x > 10 || y > 20 || x < 0 || y < 0 {
		return t.emptyColor
	}
	nx := x - t.x
	ny := y - t.y

	if nx >= 0 && nx < 4 && ny >= 0 && ny < 4 {
		block := t.current.Blocks[t.currentRotate]
		if block.On(nx, ny) == 1 {
			return t.current.Color
		}
	}
	b := t.box[y][x]
	if b == "" {
		return t.emptyColor
	}
	return b
}

func (t *Tetris) Set(x, y int, currentColor string) {
	if x > 10 || y > 20 || x < 0 || y < 0 {
		t.over = 1
		return
	}
	t.box[y][x] = currentColor
}

func (t *Tetris) On(x, y int) int {
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

func (t *Tetris) touch(block Block, x, y int) int {
	for i := 0; i != 4; i++ {
		for j := 0; j != 4; j++ {
			if block.On(i, j) == 1 && t.On(x+i, y+j) == 1 {
				return 1
			}
		}
	}
	return 0
}

func (t *Tetris) eliminate(y int) {
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
		t.setRank(t.rank)
		copy(t.box[1:y+1], t.box[:y])
		t.box[0] = [10]string{}
		for i := 0; i <= y; i++ {
			t.showRow(t.box[i], 2, 0, i)
		}
	}
}

func (t *Tetris) merge(block Block, x, y int) {
	for j := 0; j != 4; j++ {
		for i := 0; i != 4; i++ {
			if block.On(i, j) == 1 {
				t.Set(x+i, y+j, t.current.Color)
			}
		}
		t.eliminate(y + j)
	}
}

func (t *Tetris) LeftRotate() {
	currentRotate := t.currentRotate - 1
	if currentRotate < 0 {
		currentRotate = int8(len(t.current.Blocks) - 1)
	}
	if t.touch(t.current.Blocks[currentRotate], t.x, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.currentRotate = currentRotate
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	}
}

func (t *Tetris) RightRotate() {
	currentRotate := t.currentRotate + 1
	if currentRotate >= int8(len(t.current.Blocks)) {
		currentRotate = 0
	}
	if t.touch(t.current.Blocks[currentRotate], t.x, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.currentRotate = currentRotate
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	}
}

func (t *Tetris) LeftMove() {
	if t.touch(t.current.Blocks[t.currentRotate], t.x-1, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.x--
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	}
}

func (t *Tetris) RightMove() {
	if t.touch(t.current.Blocks[t.currentRotate], t.x+1, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.x++
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
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
	if t.touch(t.current.Blocks[t.currentRotate], t.x, t.y+1) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.y++
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
		return 1
	}
	t.merge(t.current.Blocks[t.currentRotate], t.x, t.y)
	t.next()
	return 0
}

func (t *Tetris) Run() (err error) {

	type Command uint

	const (
		None Command = iota
		RightRotate
		LeftRotate
		RightMove
		LeftMove
		Down
		Drop
	)

	tick := time.NewTicker(time.Second)
	cch := make(chan Command, 0)

	go func() {
		for range tick.C {
			if t.over == 1 {
				break
			}
			cch <- Down
		}
	}()

	go func() {
		for c := range cch {
			if t.over == 1 {
				return
			}
			switch c {
			case RightRotate:
				t.RightRotate()
			case LeftRotate:
				t.LeftRotate()
			case RightMove:
				t.RightMove()
			case LeftMove:
				t.LeftMove()
			case Drop:
				t.Drop()
			case Down:
				t.setTime()
				t.Down()
			}
		}
	}()

loop:
	for {
		b, _, err0 := getch.Getch()
		if err != nil {
			err = err0
			break
		}
		if t.over == 1 {
			break
		}
		c := None

		switch b {
		case 'e', 'E':
			c = RightRotate
		case 'q', 'Q':
			c = LeftRotate
		case 's', 'S':
			c = Down
		case 'a', 'A':
			c = LeftMove
		case 'd', 'D':
			c = RightMove
		case 'w', 'W':
			c = Drop
		case 'l', 'L':
			break loop
		default:
			continue
		}
		cch <- c
	}
	tick.Stop()
	t.end()
	close(cch)
	return err
}
