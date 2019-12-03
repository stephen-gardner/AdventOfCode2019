package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type intCode []int

func (code intCode) add(i int) {
	code[code[i+3]] = code[code[i+1]] + code[code[i+2]]
}

func (code intCode) multiply(i int) {
	code[code[i+3]] = code[code[i+1]] * code[code[i+2]]
}

func (code intCode) compute() {
	i := 0
	for {
		switch code[i] {
		case 99:
			return
		case 1:
			code.add(i)
		case 2:
			code.multiply(i)
		default:
			return
		}
		i += 4
	}
}

func (code intCode) print() {
	fmt.Printf("%v\n", code)
}

func (code intCode) copy() intCode {
	return append(intCode{}, code...)
}

func (code intCode) findInput(seek int) {
	fmt.Printf("Seeking %d...\n", seek)
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			state := code.copy()
			state[1] = noun
			state[2] = verb
			state.compute()
			if state[0] == seek {
				fmt.Printf("Input: %d\n", (100*noun)+verb)
				return
			}
		}
	}
}

func main() {
	data, err := ioutil.ReadFile("d02/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.Atoi(strValues[i])
	}
	state := code.copy()
	state[1] = 12
	state[2] = 2
	state.compute()
	state.print()
	code.findInput(19690720)
}
