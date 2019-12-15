package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strings"
)

type point struct {
	x, y int
}

func getAsteroids(spacemap []string) []point {
	var points []point
	for y := 0; y < len(spacemap); y++ {
		for x := 0; x < len(spacemap[0]); x++ {
			if spacemap[y][x] == '#' {
				points = append(points, point{x, y})
			}
		}
	}
	return points
}

// https://stackoverflow.com/questions/11907947/how-to-check-if-a-point-lies-on-a-line-between-2-other-points
func isOccluded(asteroids []point, origin, dest point) bool {
	for _, check := range asteroids {
		if check == origin || check == dest {
			continue
		}
		dxc := float64(check.x - origin.x)
		dyc := float64(check.y - origin.y)
		dxl := float64(dest.x - origin.x)
		dyl := float64(dest.y - origin.y)
		cross := (dxc * dyl) - (dyc * dxl)
		if cross != 0 {
			continue
		}
		occlude := false
		if math.Abs(dxl) >= math.Abs(dyl) {
			if dxl > 0 {
				occlude = origin.x <= check.x && check.x <= dest.x
			} else {
				occlude = dest.x <= check.x && check.x <= origin.x
			}
		} else {
			if dyl > 0 {
				occlude = origin.y <= check.y && check.y <= dest.y
			} else {
				occlude = dest.y <= check.y && check.y <= origin.y
			}
		}
		if occlude {
			return true
		}
	}
	return false
}

func bestAsteroid(asteroids []point) (point, int) {
	var best point
	highCount := -1
	for _, origin := range asteroids {
		count := 0
		for _, dest := range asteroids {
			if !isOccluded(asteroids, origin, dest) {
				count++
			}
		}
		if highCount == -1 || count > highCount {
			best = origin
			highCount = count
		}
	}
	return best, highCount
}

func (origin *point) distance(dest point) float64 {
	x := float64(dest.x - origin.x)
	y := float64(dest.y - origin.y)
	return math.Sqrt((x * x) + (y * y))
}

func (origin *point) radiansOffset(target point) float64 {
	rad := math.Atan2(float64(target.y-origin.y), float64(target.x-origin.x)) + (math.Pi / 2)
	if rad < 0 {
		rad += 2 * math.Pi
	}
	return rad
}

func (origin *point) createKillMap(asteroids []point) map[float64][]point {
	offsets := make(map[float64][]point, 0)
	for _, target := range asteroids {
		if target == *origin {
			continue
		}
		rad := origin.radiansOffset(target)
		if _, present := offsets[rad]; !present {
			offsets[rad] = make([]point, 0)
		}
		offsets[rad] = append(offsets[rad], target)
	}
	for _, targets := range offsets {
		sort.Slice(targets, func(i, j int) bool {
			return origin.distance(targets[i]) < origin.distance(targets[j])
		})
	}
	return offsets
}

func noScope360(radians []float64, killMap map[float64][]point, limit int) ([]float64, point, int) {
	var last point
	for i := 0; i < len(radians); i++ {
		rad := radians[i]
		hitList, present := killMap[rad]
		if !present {
			continue
		}
		last = hitList[0]
		limit--
		if limit == 0 {
			return radians, last, -1
		}
		if len(hitList)-1 == 0 {
			delete(killMap, rad)
		} else {
			killMap[rad] = hitList[1:]
		}
	}
	if len(radians) <= limit {
		return nil, last, -1
	}
	return radians, last, limit
}

func (origin *point) fireLaser(asteroids []point, limit int) point {
	killMap := origin.createKillMap(asteroids)
	radians := make([]float64, 0)
	for rad := range killMap {
		radians = append(radians, rad)
	}
	sort.Float64s(radians)
	var last point
	for limit != -1 {
		radians, last, limit = noScope360(radians, killMap, limit)
	}
	return last
}

func main() {
	data, err := ioutil.ReadFile("d10/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	spacemap := strings.Split(string(data), "\n")
	fmt.Println(strings.Join(spacemap, "\n"))
	asteroids := getAsteroids(spacemap)
	best, highCount := bestAsteroid(asteroids)
	fmt.Printf("Asteroids: %d\nX: %d\tY: %d\n\n", highCount, best.x, best.y)
	limit := 200
	last := best.fireLaser(asteroids, limit)
	fmt.Printf("Asteroid #%d Destroyed\nX: %d\tY: %d\n", limit, last.x, last.y)
	fmt.Println("Answer =", (last.x*100)+last.y)
}
