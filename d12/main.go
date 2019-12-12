package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type (
	moon struct {
		pos      vector
		velocity vector
		pulled   []bool
	}
	moonList []*moon
	vector   [3]int
)

// https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)
	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}
	return result
}

func (origin *moon) getPotentialEnergy() int {
	energy := 0
	for axis := range origin.pos {
		energy += int(math.Abs(float64(origin.pos[axis])))
	}
	return energy
}

func (origin *moon) getKineticEnergy() int {
	energy := 0
	for axis := range origin.velocity {
		energy += int(math.Abs(float64(origin.velocity[axis])))
	}
	return energy
}

func (moons moonList) applyGravity() {
	for i := range moons {
		for t := range moons {
			if moons[i] == moons[t] || moons[i].pulled[t] || moons[t].pulled[i] {
				continue
			}
			for axis := range moons[i].pos {
				if moons[i].pos[axis] < moons[t].pos[axis] {
					moons[i].velocity[axis]++
					moons[t].velocity[axis]--
				} else if moons[i].pos[axis] > moons[t].pos[axis] {
					moons[i].velocity[axis]--
					moons[t].velocity[axis]++
				}
			}
			moons[i].pulled[t] = true
			moons[t].pulled[i] = true
		}
	}
}

func (moons moonList) applyVelocity() {
	for i := range moons {
		for axis := range moons[i].velocity {
			moons[i].pos[axis] += moons[i].velocity[axis]
		}
	}
}

func (moons moonList) calcStepsToDuplicateState() int {
	instance := moons.copy()
	res := vector{}
	found := make([]bool, len(res))
	for step, nFound := 0, 0; nFound != len(res); step++ {
		instance.refresh()
		instance.applyGravity()
		instance.applyVelocity()
		for axis := range res {
			if found[axis] {
				nFound++
				continue
			}
			nFound = 0
			for i := 0; i < len(instance) && instance[i].pos[axis] == moons[i].pos[axis]; i++ {
				if i == len(instance)-1 && instance[0].velocity[axis] == 0 {
					if res[axis] == 0 {
						res[axis] = step
					} else {
						res[axis] = step - res[axis]
						found[axis] = true
					}
					break
				}
			}
		}
	}
	return LCM(res[0], res[1], res[2:]...)
}

func (moons moonList) copy() moonList {
	listCopy := make(moonList, len(moons))
	for i := range moons {
		dup := *moons[i]
		listCopy[i] = &dup
	}
	return listCopy
}

func (moons moonList) refresh() {
	for i := range moons {
		if moons[i].pulled == nil {
			moons[i].pulled = make([]bool, len(moons))
		}
		for t := range moons[i].pulled {
			moons[i].pulled[t] = false
		}
	}
}

func (moons moonList) simulate(steps int) int {
	instance := moons.copy()
	for i := 0; i < steps; i++ {
		instance.refresh()
		instance.applyGravity()
		instance.applyVelocity()
	}
	totalEnergy := 0
	for _, m := range instance {
		totalEnergy += m.getPotentialEnergy() * m.getKineticEnergy()
	}
	return totalEnergy
}

func main() {
	data, err := ioutil.ReadFile("d12/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	moons := make(moonList, len(lines))
	for i := range lines {
		moons[i] = &moon{}
		axis := strings.Split(lines[i], ",")
		for j := range axis {
			moons[i].pos[j], _ = strconv.Atoi(strings.TrimRight(strings.Split(axis[j], "=")[1], ">]"))
		}
	}
	for _, m := range moons {
		fmt.Println(m.pos)
	}
	steps := 1000
	fmt.Printf("\nTotal energy after %d steps: %d\n", steps, moons.simulate(steps))
	fmt.Printf("Steps to duplicate state: %d\n", moons.calcStepsToDuplicateState())
}
