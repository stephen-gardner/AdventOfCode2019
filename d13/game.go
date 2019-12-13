package main

import (
	"fmt"
	"os"
	"time"
)

type gameData struct {
	tiles      map[point]int64
	display    []byte
	ball       point
	paddle     point
	nextTile   point
	winWidth   int64
	winLength  int64
	writeState int
	ready      bool
}

var ticker = time.Tick(time.Millisecond * 50)
var scoreOut = point{-1, 0}

func (game *gameData) draw() {
	for pos, tileID := range game.tiles {
		if pos == scoreOut {
			continue
		}
		c := ' '
		switch tileID {
		// Wall
		case 1:
			c = '#'
		// Block
		case 2:
			c = '.'
		// Horizontal Paddle
		case 3:
			c = '='
		// Ball
		case 4:
			c = 'O'
		}
		game.display[(pos.y*game.winWidth)+pos.x] = byte(c)
	}
	score := []byte(fmt.Sprintf("%*d", game.winWidth/2, game.tiles[scoreOut]))
	for i := int64(0); i < int64(len(score)); i++ {
		game.display[((game.winLength-1)*game.winWidth)+i+1] = score[i]
	}
	_, _ = os.Stdout.Write([]byte("\x1b[2J"))
	_, _ = os.Stdout.Write(game.display)
	<-ticker
}

func (game *gameData) renderLoop() {
	if game.ready {
		game.draw()
		return
	}
	if _, present := game.tiles[scoreOut]; !present {
		return
	}
	for pos := range game.tiles {
		if pos.x+2 > game.winWidth {
			game.winWidth = pos.x + 2
		}
		if pos.y+1 > game.winLength {
			game.winLength = pos.y + 1
		}
	}
	game.display = make([]byte, game.winWidth*game.winLength)
	for y := int64(1); y <= game.winLength; y++ {
		game.display[(y*game.winWidth)-1] = '\n'
	}
	game.ready = true
}
