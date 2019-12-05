package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type (
	intCode []int

	instruction struct {
		op           int
		positionMode []bool
	}
)

func (code intCode) getInstruction(i int) instruction {
	n := code[i]
	digits := make([]int, 5)
	for i := 4; i >= 0; i-- {
		digits[i] = n % 10
		n /= 10
	}
	instr := instruction{
		op:           (digits[3] * 10) + digits[4],
		positionMode: []bool{digits[2] == 0, digits[1] == 0, digits[0] == 0},
	}
	return instr
}

func (code intCode) getVars(instr instruction, i int, forceImmediate ...bool) []int {
	vars := make([]int, len(forceImmediate))
	for n := 0; n < len(forceImmediate); n++ {
		vars[n] = code[i+n+1]
		if !forceImmediate[n] && instr.positionMode[n] {
			vars[n] = code[vars[n]]
		}
	}
	return vars
}

func (code intCode) add(instr instruction, i int) int {
	v := code.getVars(instr, i, false, false, true)
	code[v[2]] = v[0] + v[1]
	return 4
}

func (code intCode) multiply(instr instruction, i int) int {
	v := code.getVars(instr, i, false, false, true)
	code[v[2]] = v[0] * v[1]
	return 4
}

func (code intCode) input(i int) int {
	// So far only one input is possible
	code[code[i+1]] = 5
	return 2
}

func (code intCode) output(instr instruction, i int) int {
	v := code.getVars(instr, i, false)
	fmt.Println(v[0])
	return 2
}

func (code intCode) jnz(instr instruction, i int) int {
	v := code.getVars(instr, i, false, false)
	if v[0] != 0 {
		return v[1]
	}
	return i + 3
}

func (code intCode) jz(instr instruction, i int) int {
	v := code.getVars(instr, i, false, false)
	if v[0] == 0 {
		return v[1]
	}
	return i + 3
}

func (code intCode) lessThan(instr instruction, i int) int {
	v := code.getVars(instr, i, false, false, true)
	if v[0] < v[1] {
		code[v[2]] = 1
	} else {
		code[v[2]] = 0
	}
	return 4
}

func (code intCode) equals(instr instruction, i int) int {
	v := code.getVars(instr, i, false, false, true)
	if v[0] == v[1] {
		code[v[2]] = 1
	} else {
		code[v[2]] = 0
	}
	return 4
}

func (code intCode) compute() {
	i := 0
	for {
		instr := code.getInstruction(i)
		switch instr.op {
		case 1:
			i += code.add(instr, i)
		case 2:
			i += code.multiply(instr, i)
		case 3:
			i += code.input(i)
		case 4:
			i += code.output(instr, i)
		case 5:
			i = code.jnz(instr, i)
		case 6:
			i = code.jz(instr, i)
		case 7:
			i += code.lessThan(instr, i)
		case 8:
			i += code.equals(instr, i)
		case 99:
			return
		default:
			return
		}
	}
}

func (code intCode) print() {
	fmt.Printf("%v\n", code)
}

func (code intCode) copy() intCode {
	return append(intCode{}, code...)
}

func main() {
	data, err := ioutil.ReadFile("d05/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.Atoi(strValues[i])
	}
	state := code.copy()
	state.compute()
}
