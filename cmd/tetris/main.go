package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wzshiming/getch"
	"github.com/wzshiming/notify"
	"github.com/wzshiming/tetris"
)

func main() {
	err := play()
	if err != nil {
		fmt.Println(err)
	}
}

func play() error {
	var ctx, cancel = context.WithCancel(context.Background())
	notify.Once(os.Interrupt, cancel)
	cch := make(chan tetris.Command)

	go func() {
		for {
			b, _, err := getch.Getch()
			if err != nil {
				close(cch)
				fmt.Println(err)
				return
			}

			var c tetris.Command
			switch b {
			case ' ':
				c = tetris.Pause
			case 'e', 'E':
				c = tetris.RightRotate
			case 'q', 'Q':
				c = tetris.LeftRotate
			case 's', 'S':
				c = tetris.DownMove
			case 'a', 'A':
				c = tetris.LeftMove
			case 'd', 'D':
				c = tetris.RightMove
			case 'w', 'W':
				c = tetris.Drop
			case 'l', 'L':
				c = tetris.None
				cancel()
			default:
				c = tetris.None
			}
			cch <- c
		}
	}()

	t := tetris.NewTetris()
	return t.Run(ctx, cch)
}
