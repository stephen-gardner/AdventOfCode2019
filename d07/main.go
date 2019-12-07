package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

type phases []int
type amplifiers []*process

func (ph phases) toInt() int {
	n := 0
	for i := 0; i < len(ph); i++ {
		n = (n * 10) + ph[i]
	}
	return n
}

func (ph phases) validPhase() bool {
	for i := 0; i < len(ph); i++ {
		if ph[i] < 5 {
			return false
		}
		for j := i + 1; j < len(ph); j++ {
			if ph[i] == ph[j] {
				return false
			}
		}
	}
	return true
}

func (ph phases) getNext() phases {
	for {
		tmp := ph.toInt() + 1
		if tmp > 98765 {
			return nil
		}
		for i := 4; i >= 0; i-- {
			ph[i] = tmp % 10
			tmp /= 10
		}
		if ph.validPhase() {
			return ph
		}
	}
}

func (amps amplifiers) initAmplifiers(code intCode) {
	for i := 0; i < len(amps); i++ {
		amps[i] = &process{
			code:   code.copy(),
			stdout: make(chan int, 2),
		}
	}
	for i := 0; i < len(amps); i++ {
		amps[(i+1)%len(amps)].stdin = amps[i].stdout
	}
}

func (amps amplifiers) run(phaseSettings phases) {
	var wg sync.WaitGroup
	wg.Add(len(amps))
	for i := 0; i < len(amps); i++ {
		args := []int{phaseSettings[i]}
		if i == 0 {
			args = append(args, 0)
		}
		go amps[i].compute(&wg, args)
	}
	wg.Wait()
}

func findOptimalPhase(code intCode) {
	var optimalPhases phases
	phaseSettings := phases{5, 6, 7, 8, 9}
	maxSignal := 0
	for {
		amplifiers := make(amplifiers, 5)
		amplifiers.initAmplifiers(code)
		amplifiers.run(phaseSettings)
		signal := <-amplifiers[4].stdout
		if signal > maxSignal {
			optimalPhases = append([]int{}, phaseSettings...)
			maxSignal = signal
		}
		phaseSettings = phaseSettings.getNext()
		if phaseSettings == nil {
			break
		}
	}
	fmt.Printf("Optimal Phase: %v\nHighest signal: %d\n", optimalPhases, maxSignal)
}

func main() {
	data, err := ioutil.ReadFile("d07/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.Atoi(strValues[i])
	}
	findOptimalPhase(code)
}
