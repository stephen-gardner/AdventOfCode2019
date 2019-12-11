package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type (
	point struct {
		x, y int
	}

	location struct {
		coords point
		face   int
	}
)

func (proc *process) paint(paintMap map[point]int64, startingColor int64) {
	pos := location{}
	paint := true
	fmt.Println("Starting painting robot with arg:", startingColor)
	go proc.compute([]int64{startingColor})
	for out := range proc.stdout {
		if paint {
			paintMap[pos.coords] = out
			paint = false
			continue
		}
		dir := 1
		if out == 0 {
			dir = -1
		}
		switch pos.face {
		case 0:
			pos.coords.x += dir
		case 1:
			pos.coords.y += dir
		case 2:
			pos.coords.x -= dir
		case 3:
			pos.coords.y -= dir
		}
		pos.face = (pos.face - dir + 4) % 4
		proc.stdin <- paintMap[pos.coords]
		paint = true
	}
}

func findBounds(paintMap map[point]int64, rendered int64) (int, int, int, int) {
	var offX, offY, highX, highY int
	first := true
	for pos, color := range paintMap {
		if color != rendered {
			continue
		}
		if first {
			offX, offY, highX, highY = pos.x, pos.y, pos.x, pos.y
			first = false
		}
		if pos.x < offX {
			offX = pos.x
		}
		if pos.x > highX {
			highX = pos.x
		}
		if pos.y < offY {
			offY = pos.y
		}
		if pos.y > highY {
			highY = pos.y
		}
	}
	offX *= -1
	offY *= -1
	highX += offX
	highY += offY
	return offX, offY, highX, highY
}

func printPaintMap(paintMap map[point]int64) {
	offX, offY, highX, highY := findBounds(paintMap, 1)
	output := make([][]byte, highY+1)
	for i := range output {
		output[i] = []byte(strings.Repeat(".", highX+1))
	}
	for pos, color := range paintMap {
		if color == 1 {
			output[len(output)-(pos.y+offY+1)][pos.x+offX] = '#'
		}
	}
	for _, line := range output {
		fmt.Println(strings.ReplaceAll(string(line), ".", " "))
	}
}

func runPaintingBot(code intCode, startingColor int64) map[point]int64 {
	robot := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64, 1),
		stdout: make(chan int64),
	}
	paintMap := make(map[point]int64)
	robot.paint(paintMap, startingColor)
	return paintMap
}

func main() {
	data, err := ioutil.ReadFile("d11/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	fmt.Printf("Total painted panels: %d\n\n", len(runPaintingBot(code, 0)))
	printPaintMap(runPaintingBot(code, 1))
}
