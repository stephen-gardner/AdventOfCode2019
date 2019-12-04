package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type (
	coords struct {
		x, y int
	}
	wire struct {
		dir []string
		pos coords
	}
)

var routing = make(map[coords]int)
var intersections = make(map[coords]bool)

func (pos coords) route(steps int, first bool) {
	if first {
		if _, overlap := routing[pos]; !overlap {
			routing[pos] = steps
		}
	} else {
		if _, overlap := routing[pos]; overlap {
			if _, exists := intersections[pos]; !exists {
				routing[pos] += steps
				intersections[pos] = true
			}
		}
	}
}

func (wire *wire) up(n, steps int, first bool) {
	for ; n > 0; n-- {
		wire.pos.y++
		steps++
		wire.pos.route(steps, first)
	}
}

func (wire *wire) down(n, steps int, first bool) {
	for ; n > 0; n-- {
		wire.pos.y--
		steps++
		wire.pos.route(steps, first)
	}
}

func (wire *wire) left(n, steps int, first bool) {
	for ; n > 0; n-- {
		wire.pos.x--
		steps++
		wire.pos.route(steps, first)
	}
}

func (wire *wire) right(n, steps int, first bool) {
	for ; n > 0; n-- {
		wire.pos.x++
		steps++
		wire.pos.route(steps, first)
	}
}

func (wire *wire) route(first bool) {
	steps := 0
	for range wire.dir {
		dir := wire.dir[0][0]
		mag, _ := strconv.Atoi(wire.dir[0][1:])
		switch dir {
		case 'U':
			wire.up(mag, steps, first)
		case 'D':
			wire.down(mag, steps, first)
		case 'L':
			wire.left(mag, steps, first)
		case 'R':
			wire.right(mag, steps, first)
		default:
			break
		}
		wire.dir = wire.dir[1:]
		steps += mag
	}
}

func findClosestIntersection() int {
	closest := -1
	for pos := range intersections {
		mDist := int(math.Abs(float64(pos.x)) + math.Abs(float64(pos.y)))
		if closest == -1 || mDist < closest {
			closest = mDist
		}
	}
	return closest
}

func findShortestIntersection() int {
	shortest := -1
	for pos := range intersections {
		steps := routing[pos]
		if shortest == -1 || steps < shortest {
			shortest = steps
		}
	}
	return shortest
}

func main() {
	file, err := os.Open("d03/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	wireData1, _ := reader.ReadString('\n')
	wireData2, _ := reader.ReadString('\n')
	w1 := wire{
		dir: strings.Split(wireData1, ","),
		pos: coords{0, 0},
	}
	w2 := wire{
		dir: strings.Split(wireData2, ","),
		pos: coords{0, 0},
	}
	w1.route(true)
	w2.route(false)
	if len(intersections) > 0 {
		fmt.Printf("Closest intersection has a manhatten distance of: %d\n", findClosestIntersection())
		fmt.Printf("Shortest intersection: %d\n", findShortestIntersection())
	}
}
