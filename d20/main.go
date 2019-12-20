package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

type point struct {
	x int
	y int
}

type position struct {
	coords point
	steps  int
	level  int
}

const NORTH int = 0
const SOUTH int = 1
const WEST int = 2
const EAST int = 3

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

func findPortals(data [][]byte, portals map[point]point) (start point, end point) {
	portalIndex := make(map[string]point)
	for row := range data {
		for col := range data[row] {
			pos := point{col, row}
			if c := data[row][col]; c >= 'A' && c <= 'Z' {
				for dir := range dirs {
					adj := pos.add(dir)
					if adj.y < 0 || adj.x < 0 || adj.y > len(data)-1 || adj.x > len(data[adj.y])-1 {
						continue
					}
					c2 := data[adj.y][adj.x]
					if !(c2 >= 'A' && c2 <= 'Z') {
						continue
					}
					portal := adj.add(dir)
					if portal.y > len(data)-1 || portal.x > len(data[portal.y])-1 || data[portal.y][portal.x] != '.' {
						if dir == NORTH {
							portal = pos.add(SOUTH)
						} else {
							portal = pos.add(WEST)
						}
					}
					label := string([]byte{c, c2})
					if dest, present := portalIndex[label]; present {
						portals[portal] = dest
						portals[dest] = portal
					} else {
						portalIndex[label] = portal
					}
					if label == "AA" {
						start = portal
					}
					if label == "ZZ" {
						end = portal
					}
					data[pos.y][pos.x] = ' '
					data[adj.y][adj.x] = ' '
					break
				}
			}
		}
	}
	return
}

func makeRecursive(data [][]byte, portals map[point]point) {
	for portal := range portals {
		if portal.y == 2 || portal.y == len(data)-3 || portal.x == 2 || portal.x == len(data[portal.y])-3 {
			data[portal.y][portal.x] = 'O'
		} else {
			data[portal.y][portal.x] = 'I'
		}
	}
}

func BFS(data [][]byte, portals map[point]point, start, end point) int {
	queue := []position{{start, 0, 0}}
	visited := make([]map[point]bool, 1)
	visited[0] = make(map[point]bool)
	visited[0][start] = true
	for len(queue) > 0 {
		for _, origin := range queue {
			queue = queue[1:]
			for dir := range dirs {
				pos := position{origin.coords.add(dir), origin.steps + 1, origin.level}
				if pos.coords.y < 0 || pos.coords.x < 0 ||
					pos.coords.y > len(data)-1 || pos.coords.x > len(data[pos.coords.y])-1 {
					continue
				}
				if _, present := visited[pos.level][pos.coords]; present {
					continue
				}
				c := data[pos.coords.y][pos.coords.x]
				if c != '.' && c != 'I' && c != 'O' {
					continue
				}
				visited[pos.level][pos.coords] = true
				if pos.coords == end && pos.level == 0 {
					return pos.steps
				}
				if dest, isPortal := portals[pos.coords]; isPortal {
					if c == 'I' {
						pos.level++
					} else if c == 'O' {
						pos.level--
					}
					if pos.level < 0 || pos.level > 25 {
						continue
					}
					if pos.level > len(visited)-1 {
						visited = append(visited, make(map[point]bool))
					}
					pos.coords = dest
					visited[pos.level][pos.coords] = true
					pos.steps += 1
				}
				queue = append(queue, pos)
			}
		}
	}
	return -1
}

func main() {
	input, err := ioutil.ReadFile("d20/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.Split(input, []byte("\n"))
	for _, line := range data {
		fmt.Println(string(line))
	}
	portals := make(map[point]point)
	start, end := findPortals(data, portals)
	fmt.Println("Steps to exit:", BFS(data, portals, start, end))
	makeRecursive(data, portals)
	fmt.Println("Steps to exit (recursive):", BFS(data, portals, start, end))
}
