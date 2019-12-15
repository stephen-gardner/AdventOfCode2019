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
		x int64
		y int64
	}
	location struct {
		coords point
		dir    int64
	}
)

const NORTH int64 = 1
const SOUTH int64 = 2
const WEST int64 = 3
const EAST int64 = 4

const WALL int64 = 0
const PATH int64 = 1
const SYSTEM int64 = 2
const DROID int64 = 3
const OXYGEN int64 = 42
const EXPLORE int64 = -1

var dirs = []point{
	{0, 1},
	{0, -1},
	{-1, 0},
	{1, 0},
}

func (p point) add(dir int64) point {
	p2 := dirs[dir-1]
	p.x += p2.x
	p.y += p2.y
	return p
}

func (proc *process) droidDFS(pos location, area map[point]int64, target int64) (int64, point) {
	if _, visited := area[pos.coords]; visited {
		return 0, pos.coords
	}
	proc.stdin <- pos.dir
	status := <-proc.stdout
	area[pos.coords] = status
	if status == WALL {
		return 0, pos.coords
	} else if status == target {
		return 1, pos.coords
	}
	var dest point
	shortest := int64(-1)
	for i := int64(1); i < 5; i++ {
		if steps, pos := proc.droidDFS(location{pos.coords.add(i), i}, area, target); steps > 0 {
			if shortest == -1 || steps < shortest {
				shortest = steps
				dest = pos
			}
		}
		if shortest != -1 {
			return 1 + shortest, dest
		}
	}
	switch pos.dir {
	case NORTH:
		proc.stdin <- SOUTH
	case SOUTH:
		proc.stdin <- NORTH
	case EAST:
		proc.stdin <- WEST
	case WEST:
		proc.stdin <- EAST
	default:
		return 0, pos.coords
	}
	<-proc.stdout
	area[pos.coords] = EXPLORE
	return 0, pos.coords
}

func findBounds(areaMap map[point]int64) (int64, int64, int64, int64) {
	var offX, offY, highX, highY int64
	first := true
	for pos := range areaMap {
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

func printAreaMap(areaMap map[point]int64) {
	offX, offY, highX, highY := findBounds(areaMap)
	output := make([][]byte, highY+1)
	for i := range output {
		output[i] = []byte(strings.Repeat(" ", int(highX)+1))
	}
	for pos, status := range areaMap {
		var c byte
		switch status {
		case PATH:
			c = '.'
		case WALL:
			c = '#'
		case SYSTEM:
			c = 'T'
		case DROID:
			c = 'D'
		case OXYGEN:
			c = 'O'
		default:
			c = ' '
		}
		output[int64(len(output))-(pos.y+offY+1)][pos.x+offX] = c
	}
	for _, line := range output {
		fmt.Println(string(line))
	}
}

func runDroid(code intCode, target int64) (int64, map[point]int64, point) {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	area := make(map[point]int64)
	go proc.compute([]int64{})
	n, targetPos := proc.droidDFS(location{dir: NORTH}, area, target)
	area[point{0, 0}] = 3
	return n, area, targetPos
}

func oxygenate(area map[point]int64, source point) int {
	minutes := 0
	queue := []point{source}
	for len(queue) > 0 {
		for _, pos := range queue {
			queue = queue[1:]
			if area[pos] != OXYGEN {
				continue
			}
			for i := int64(1); i < 5; i++ {
				newPos := pos.add(i)
				if status, visited := area[newPos]; visited && (status == EXPLORE || status == PATH) {
					area[newPos] = OXYGEN
					queue = append(queue, newPos)
				}
			}
		}
		if len(queue) > 0 {
			minutes++
		}
	}
	return minutes
}

func main() {
	data, err := ioutil.ReadFile("d15/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	steps, solution, systemPos := runDroid(code, SYSTEM)
	_, area, _ := runDroid(code, EXPLORE)
	for pos, status := range solution {
		area[pos] = status
	}
	printAreaMap(area)
	area[point{0, 0}] = PATH
	area[systemPos] = OXYGEN
	minutes := oxygenate(area, systemPos)
	fmt.Printf("Reached target in %d steps\n", steps)
	fmt.Printf("Filled area with oxygen in %d minutes\n", minutes)
}
