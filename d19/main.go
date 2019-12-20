package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func beam(code intCode, x, y int) bool {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	go proc.compute([]int64{int64(x), int64(y)})
	return <-proc.stdout != 0
}

func findSquare(code intCode, size, xMax, yMax int) (point, bool) {
	xOffset := 0
	size--
	for y := size; y < yMax; y++ {
		for x := xOffset; x < xMax; x++ {
			if beam(code, x, y) {
				if x > xOffset {
					xOffset = x
				}
				if beam(code, x, y-size) && beam(code, x+size, y-size) {
					return point{x, y - size}, true
				}
				x = xMax
			}
		}
	}
	return point{}, false
}

func findCoverage(code intCode, xMax, yMax int) int {
	xOffset := 0
	hits := 0
	for y := 0; y < yMax; y++ {
		rowStart := true
		for x := xOffset; x < xMax; x++ {
			if beam(code, x, y) {
				if rowStart {
					rowStart = false
					if x > xOffset {
						xOffset = x
					}
				}
				hits++
			} else if !rowStart {
				x = xMax
			}
		}
	}
	return hits
}

func main() {
	data, err := ioutil.ReadFile("d19/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	fmt.Println("Tractor beam coverage for 50x50 area:", findCoverage(code, 50, 50))
	res, _ := findSquare(code, 100, 5000, 5000)
	fmt.Println("Scan for 100x100 square:", (res.x*10000)+res.y)
}
