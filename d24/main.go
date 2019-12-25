package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type point struct {
	x int
	y int
}

var dirs = []point{
	{0, 1},
	{0, -1},
	{-1, 0},
	{1, 0},
}

func (p point) add(dir int) point {
	p2 := dirs[dir]
	p.x += p2.x
	p.y += p2.y
	return p
}

func sumCol(data [][]byte, col int) int {
	bugs := 0
	for row := 0; row < len(data); row++ {
		if data[row][col] == '#' {
			bugs++
		}
	}
	return bugs
}

func sumRow(data [][]byte, row int) int {
	bugs := 0
	for col := 0; col < len(data[row]); col++ {
		if data[row][col] == '#' {
			bugs++
		}
	}
	return bugs
}

func isBlank(data [][]byte) bool {
	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[row]); col++ {
			if data[row][col] == '#' {
				return false
			}
		}
	}
	return true
}

func getNextRecursivePhase(levels map[int][][]byte) {
	res := make(map[int][][]byte)
	highLevel := 1
	lowLevel := -1
	for level := range levels {
		if level < lowLevel {
			lowLevel = level
		}
		if level > highLevel {
			highLevel = level
		}
	}
	if !isBlank(levels[highLevel]) {
		levels[highLevel+1] = getBlankSlate(levels[highLevel])
		highLevel++
	}
	if !isBlank(levels[lowLevel]) {
		levels[lowLevel-1] = getBlankSlate(levels[lowLevel])
		lowLevel--
	}
	for level, data := range levels {
		slate := getBlankSlate(data)
		for row := 0; row < len(data); row++ {
			for col := 0; col < len(data[row]); col++ {
				if row == 2 && col == 2 {
					continue
				}
				adjBugs := 0
				origin := point{col, row}
				for dir := range dirs {
					pos := origin.add(dir)
					if level == lowLevel &&
						(pos.y < 0 || pos.y >= len(data) || pos.x < 0 || pos.x >= len(data[pos.y])) {
						continue
					}
					if level == highLevel && pos.x == 2 && pos.y == 2 {
						continue
					}
					var ref byte
					switch {
					case pos.y < 0:
						ref = levels[level-1][1][2]
					case pos.y >= len(data):
						ref = levels[level-1][3][2]
					case pos.x < 0:
						ref = levels[level-1][2][1]
					case pos.x >= len(data[pos.y]):
						ref = levels[level-1][2][3]
					case pos.y == 2 && pos.x == 2:
						switch {
						case col == 1:
							adjBugs += sumCol(levels[level+1], 0)
						case col == 3:
							adjBugs += sumCol(levels[level+1], 4)
						case row == 1:
							adjBugs += sumRow(levels[level+1], 0)
						case row == 3:
							adjBugs += sumRow(levels[level+1], 4)
						}
						continue
					default:
						ref = data[pos.y][pos.x]
					}
					if ref == '#' {
						adjBugs++
					}
				}
				if data[row][col] == '#' {
					if adjBugs == 1 {
						slate[row][col] = '#'
					}
				} else if adjBugs == 1 || adjBugs == 2 {
					slate[row][col] = '#'
				}
			}
		}
		res[level] = slate
	}
	for level, data := range res {
		levels[level] = data
	}
}

func countBugs(rData map[int][][]byte) int {
	bugs := 0
	for _, data := range rData {
		for row := 0; row < len(data); row++ {
			for col := 0; col < len(data[row]); col++ {
				if row == 2 && col == 2 {
					continue
				}
				if data[row][col] == '#' {
					bugs++
				}
			}
		}
	}
	return bugs
}

func recursiveBugs(data [][]byte, nPhases int) {
	levels := make(map[int][][]byte)
	levels[-1] = getBlankSlate(data)
	levels[0] = data
	levels[1] = getBlankSlate(data)
	for i := 0; i < nPhases; i++ {
		getNextRecursivePhase(levels)
	}
	fmt.Println("Total bugs (recursive):", countBugs(levels))
}

func getNextPhase(data [][]byte) [][]byte {
	res := make([][]byte, len(data))
	for i := range data {
		res[i] = []byte(strings.Repeat(".", len(data[i])))
	}
	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[row]); col++ {
			origin := point{col, row}
			adjBugs := 0
			for dir := range dirs {
				pos := origin.add(dir)
				if pos.y < 0 || pos.y >= len(data) || pos.x < 0 || pos.x >= len(data[pos.y]) {
					continue
				}
				if data[pos.y][pos.x] == '#' {
					adjBugs++
				}
			}
			if data[row][col] == '#' {
				if adjBugs == 1 {
					res[row][col] = '#'
				}
			} else if adjBugs == 1 || adjBugs == 2 {
				res[row][col] = '#'
			}
		}
	}
	return res
}

func calculateRating(data [][]byte) int {
	res := 0
	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[row]); col++ {
			if data[row][col] == '#' {
				res |= 1 << uint((row*len(data[0]))+col)
			}
		}
	}
	return res
}

func getBlankSlate(data [][]byte) [][]byte {
	blank := make([][]byte, len(data))
	for i := range data {
		blank[i] = []byte(strings.Repeat(".", len(data[i])))
	}
	return blank
}

func findDuplicateState(originalData [][]byte) {
	data := getBlankSlate(originalData)
	for row := 0; row < len(data); row++ {
		for col := 0; col < len(data[row]); col++ {
			data[row][col] = originalData[row][col]
		}
	}
	states := make(map[int]bool)
	states[calculateRating(data)] = true
	for {
		data = getNextPhase(data)
		rating := calculateRating(data)
		if _, present := states[rating]; present {
			fmt.Println("Biodiversity rating:", rating)
			break
		}
		states[rating] = true
	}
}

func main() {
	input, err := ioutil.ReadFile("d24/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.Split(input, []byte("\n"))
	fmt.Println(string(bytes.Join(data, []byte("\n"))))
	findDuplicateState(data)
	recursiveBugs(data, 200)
}
