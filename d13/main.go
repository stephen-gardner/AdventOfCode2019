package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func runGame(code intCode, freePlay bool) {
	proc := &process{
		game: &gameData{tiles: make(map[point]int64)},
		code: append(code.copy(), make(intCode, 4096)...),
	}
	if freePlay {
		proc.code[0] = 2
	}
	proc.compute([]int64{})
	if !freePlay {
		blocks := 0
		for _, tileID := range proc.game.tiles {
			if tileID == 2 {
				blocks++
			}
		}
		fmt.Println("Total Blocks destroyed:", blocks)
	}
}

func main() {
	data, err := ioutil.ReadFile("d13/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	runGame(code, true)
	runGame(code, false)
}
