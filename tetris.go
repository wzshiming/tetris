package tetris

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Tetris struct {
	offX, offY    int
	interval      time.Duration
	start         time.Time
	pauseDuration time.Duration
	pauseTime     time.Time
	draw          *Draw
	box           [20][10]string
	x, y          int
	waiting       Blocks
	currentRotate int8
	current       Blocks
	emptyColor    string
	rand          *rand.Rand
	rank          uint64
	flag          flag
}

func NewTetris() *Tetris {
	t := &Tetris{
		interval: time.Second,
		start:    time.Now(),
		draw:     NewDraw(os.Stdout),
		rand:     rand.New(rand.NewSource(time.Now().Unix())),
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
	helpY := 14
	t.draw.Dot("  H E L P", 2, helpX, helpY)
	t.draw.Dot("L:     Quit", 2, helpX, helpY+1)
	t.draw.Dot("Q:     Left rotate", 2, helpX, helpY+2)
	t.draw.Dot("E:     Right rotate", 2, helpX, helpY+3)
	t.draw.Dot("A:     Left move", 2, helpX, helpY+4)
	t.draw.Dot("D:     Right move", 2, helpX, helpY+5)
	t.draw.Dot("S:     Down move", 2, helpX, helpY+6)
	t.draw.Dot("W:     Drop", 2, helpX, helpY+7)
	t.draw.Dot("Space: Pause", 2, helpX, helpY+8)
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
	t.draw.Dot(fmt.Sprintf("Rank:  %d", i), 2, 14, 10)
}

func (t *Tetris) setTime() {
	t.draw.Dot(fmt.Sprintf("Time:  %v\t", (time.Now().Sub(t.start)-t.pauseDuration)/time.Second*time.Second), 2, 14, 11)
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
	t.predict(false)
	t.waiting = BlocksPool[t.rand.Int()%len(BlocksPool)]
	t.predict(true)
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

func (t *Tetris) eliminate(y int) bool {
	if y >= 20 || y < 0 {
		return false
	}

	for _, d := range t.box[y] {
		if d == "" {
			return false
		}
	}

	t.rank++
	t.setRank(t.rank)
	copy(t.box[1:y+1], t.box[:y])
	t.box[0] = [10]string{}
	return true
}

func (t *Tetris) check(block Block, x, y int) bool {
	for j := 0; j != 4; j++ {
		for i := 0; i != 4; i++ {
			if block.On(i, j) == 1 {
				x := x + i
				y := y + j
				if x > 10 || y > 20 || x < 0 || y < 0 {
					return false
				}
			}
		}
	}
	return true
}

func (t *Tetris) merge(block Block, x, y int) {
	flash := false
	for j := 0; j != 4; j++ {
		for i := 0; i != 4; i++ {
			if block.On(i, j) == 1 {
				t.Set(x+i, y+j, t.current.Color)
			}
		}
		flash = t.eliminate(y+j) || flash
	}

	if flash {
		j := y + 4
		if j >= 20 {
			j = 19
		}
		for i := 0; i <= j; i++ {
			t.showRow(t.box[i], 2, 0, i)
		}
	}
}

func (t *Tetris) LeftRotate() {
	currentRotate := t.currentRotate - 1
	if currentRotate < 0 {
		currentRotate = int8(len(t.current.Blocks) - 1)
	}
	if t.touch(t.current.Blocks[currentRotate], t.x, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.predict(false)
		t.currentRotate = currentRotate
		t.predict(true)
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
		t.predict(false)
		t.currentRotate = currentRotate
		t.predict(true)
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	}
}

func (t *Tetris) LeftMove() {
	if t.touch(t.current.Blocks[t.currentRotate], t.x-1, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.predict(false)
		t.x--
		t.predict(true)
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	}
}

func (t *Tetris) RightMove() {
	if t.touch(t.current.Blocks[t.currentRotate], t.x+1, t.y) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.predict(false)
		t.x++
		t.predict(true)
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	}
}

func (t *Tetris) Drop() {
	i := 1
	for t.touch(t.current.Blocks[t.currentRotate], t.x, t.y+i) == 0 {
		i++
	}
	i -= 1

	t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
	t.y += i
	t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	if !t.check(t.current.Blocks[t.currentRotate], t.x, t.y) {
		t.flag.On(overFlag)
		return
	}
	t.merge(t.current.Blocks[t.currentRotate], t.x, t.y)
	t.next()
}

func (t *Tetris) DownMove() {
	if t.touch(t.current.Blocks[t.currentRotate], t.x, t.y+1) == 0 {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y)
		t.y++
		t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
		return
	}
	t.showBlock(t.current.Blocks[t.currentRotate], t.current.Color, 2, t.x, t.y)
	if !t.check(t.current.Blocks[t.currentRotate], t.x, t.y) {
		t.flag.On(overFlag)
		return
	}
	t.merge(t.current.Blocks[t.currentRotate], t.x, t.y)
	t.next()
}

func (t *Tetris) Pause() {
	t.flag.On(pauseFlag)
	t.pauseTime = time.Now()
}

func (t *Tetris) Continue() {
	t.flag.Off(pauseFlag)
	t.pauseDuration += time.Now().Sub(t.pauseTime)
}

func (t *Tetris) IsPause() bool {
	return t.flag.Has(pauseFlag)
}

func (t *Tetris) predict(show bool) {
	i := 1
	for t.touch(t.current.Blocks[t.currentRotate], t.x, t.y+i) == 0 {
		i++
	}
	i -= 1

	if show {
		if i != 0 {
			t.showBlock(t.current.Blocks[t.currentRotate], t.current.Predict, 2, t.x, t.y+i)
		}
	} else {
		t.showBlock(t.current.Blocks[t.currentRotate], t.emptyColor, 2, t.x, t.y+i)
	}
}

func (t *Tetris) IsGameOver() bool {
	return t.flag.Has(overFlag)
}

func (t *Tetris) Run(ctx context.Context, cch <-chan Command) (err error) {
	ticker := time.NewTimer(t.interval)
loop:
	for {
		select {
		case <-ctx.Done():
			t.flag.On(overFlag)
			break loop
		case <-ticker.C:
			if t.IsGameOver() {
				break loop
			}
			if t.IsPause() {
				continue
			}
			t.setTime()
			t.DownMove()
			ticker.Reset(t.interval)
		case c, ok := <-cch:
			if c == None {
				continue
			}
			if !ok || t.IsGameOver() {
				break loop
			}
			if t.IsPause() {
				t.Continue()
				ticker.Reset(t.interval)
				if c == Pause {
					continue
				}
			}
			switch c {
			case Pause:
				t.Pause()
			case RightRotate:
				t.RightRotate()
			case LeftRotate:
				t.LeftRotate()
			case RightMove:
				t.RightMove()
			case LeftMove:
				t.LeftMove()
			case DownMove:
				t.setTime()
				t.DownMove()
				ticker.Reset(t.interval)
			case Drop:
				t.Drop()
			}
			if t.IsGameOver() {
				break loop
			}
		}
	}

	t.end()
	return err
}
