package tetris

import (
	"testing"
)

func TestBlock(t *testing.T) {
	for _, blocks := range BlocksPool {
		for _, block := range blocks {
			t.Log("\n" + string(block.Draw("1", "0", "\n")))
		}
	}
}
