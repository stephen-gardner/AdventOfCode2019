package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

type (
	point struct {
		x int
		y int
	}
	position struct {
		coords point
		keys   int
		steps  int
	}
)

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

func findStartPositions(data [][]byte) []position {
	var positions []position
	for row := range data {
		for col := range data[row] {
			if data[row][col] == '@' {
				positions = append(positions, position{point{col, row}, 0, 0})
			}
		}
	}
	return positions
}

func updateStartPositions(data [][]byte) {
	for row := range data {
		for col := range data[row] {
			if data[row][col] == '@' {
				data[row][col] = '#'
				data[row+1][col] = '#'
				data[row-1][col] = '#'
				data[row][col+1] = '#'
				data[row][col-1] = '#'
				data[row-1][col-1] = '@'
				data[row-1][col+1] = '@'
				data[row+1][col-1] = '@'
				data[row+1][col+1] = '@'
				return
			}
		}
	}
}

func findAllKeys(data [][]byte) int {
	keys := 0
	for row := range data {
		for _, c := range data[row] {
			if c >= 'a' && c <= 'z' {
				keys |= 1 << (c - 'a')
			}
		}
	}
	return keys
}

func findAvailableKeys(data [][]byte, start point) int {
	keys := 0
	queue := []point{start}
	visited := make(map[point]bool)
	for len(queue) > 0 {
		for _, origin := range queue {
			queue = queue[1:]
			for dir := 0; dir < len(dirs); dir++ {
				pos := origin.add(dir)
				if pos.y < 0 || pos.y >= len(data) || pos.x < 0 || pos.x >= len(data[pos.y]) {
					continue
				}
				if _, seen := visited[pos]; seen {
					continue
				}
				visited[pos] = true
				c := data[pos.y][pos.x]
				if c == '#' {
					continue
				}
				if c >= 'a' && c <= 'z' {
					keys |= 1 << (c - 'a')
				}
				queue = append(queue, pos)
			}
		}
	}
	return keys
}

func findMinimalSteps(data [][]byte, start position, allKeys int, ignoreOtherDoors bool) int {
	queue := []position{start}
	visited := make(map[int]map[point]bool)
	visited[start.keys] = make(map[point]bool)
	for len(queue) > 0 {
		for _, origin := range queue {
			queue = queue[1:]
			for dir := 0; dir < len(dirs); dir++ {
				pos := origin.coords.add(dir)
				keys := origin.keys
				if pos.y < 0 || pos.y >= len(data) || pos.x < 0 || pos.x >= len(data[pos.y]) {
					continue
				}
				if _, hasState := visited[keys]; hasState {
					if _, seen := visited[keys][pos]; seen {
						continue
					}
				}
				visited[keys][pos] = true
				c := data[pos.y][pos.x]
				if c == '#' {
					continue
				}
				if c >= 'A' && c <= 'Z' {
					if keys&(1<<(c-'A')) == 0 {
						if !ignoreOtherDoors || allKeys&(1<<(c-'A')) == 1 {
							continue
						}
					}
				}
				if c >= 'a' && c <= 'z' {
					keys |= 1 << (c - 'a')
					if keys == allKeys {
						return origin.steps + 1
					}
					if _, present := visited[keys]; !present {
						// Need fresh map to cover old ground
						visited[keys] = make(map[point]bool)
						visited[keys][pos] = true
					}
				}
				queue = append(queue, position{pos, keys, origin.steps + 1})
			}
		}
	}
	return -1
}

func main() {
	input, err := ioutil.ReadFile("d18/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(input))
	data := bytes.Split(input, []byte("\n"))
	fmt.Println(
		"Steps for shortest path:",
		findMinimalSteps(data, findStartPositions(data)[0], findAllKeys(data), false),
	)
	steps := 0
	updateStartPositions(data)
	for _, pos := range findStartPositions(data) {
		steps += findMinimalSteps(data, pos, findAvailableKeys(data, pos.coords), true)
	}
	fmt.Println("Steps for shortest path (multiple bots):", steps)
}
