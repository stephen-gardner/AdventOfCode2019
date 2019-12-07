package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type (
	intCode []int

	process struct {
		code   intCode
		stdin  chan int
		stdout chan int
		pc     int
		instr  instruction
	}

	instruction struct {
		op           int
		positionMode []bool
	}
)

func (code intCode) print() {
	fmt.Printf("%v\n", code)
}

func (code intCode) copy() intCode {
	return append(intCode{}, code...)
}

// Opcode #1
func (proc *process) add() {
	v := proc.getOperands(false, false, true)
	proc.code[v[2]] = v[0] + v[1]
	proc.pc += 4
}

// Opcode #2
func (proc *process) mul() {
	v := proc.getOperands(false, false, true)
	proc.code[v[2]] = v[0] * v[1]
	proc.pc += 4
}

// Opcode #3
func (proc *process) read(args *[]int) {
	if len(*args) > 0 {
		proc.code[proc.code[proc.pc+1]] = (*args)[0]
		*args = (*args)[1:]
	} else {
		proc.code[proc.code[proc.pc+1]] = <-proc.stdin
	}
	proc.pc += 2
}

// Opcode #4
func (proc *process) write() {
	v := proc.getOperands(false)
	proc.stdout <- v[0]
	proc.pc += 2
}

// Opcode #5
func (proc *process) jnz() {
	v := proc.getOperands(false, false)
	if v[0] != 0 {
		proc.pc = v[1]
		return
	}
	proc.pc += 3
}

// Opcode #6
func (proc *process) jz() {
	v := proc.getOperands(false, false)
	if v[0] == 0 {
		proc.pc = v[1]
		return
	}
	proc.pc += 3
}

// Opcode #7
func (proc *process) lessThan() {
	v := proc.getOperands(false, false, true)
	if v[0] < v[1] {
		proc.code[v[2]] = 1
	} else {
		proc.code[v[2]] = 0
	}
	proc.pc += 4
}

// Opcode #8
func (proc *process) equals() {
	v := proc.getOperands(false, false, true)
	if v[0] == v[1] {
		proc.code[v[2]] = 1
	} else {
		proc.code[v[2]] = 0
	}
	proc.pc += 4
}

func (proc *process) getOperands(forceImmediate ...bool) []int {
	vars := make([]int, len(forceImmediate))
	for n := 0; n < len(forceImmediate); n++ {
		vars[n] = proc.code[proc.pc+n+1]
		if !forceImmediate[n] && proc.instr.positionMode[n] {
			vars[n] = proc.code[vars[n]]
		}
	}
	return vars
}

func (proc *process) decodeInstruction() {
	n := proc.code[proc.pc]
	digits := make([]int, 5)
	for i := 4; i >= 0; i-- {
		digits[i] = n % 10
		n /= 10
	}
	proc.instr = instruction{
		op:           (digits[3] * 10) + digits[4],
		positionMode: []bool{digits[2] == 0, digits[1] == 0, digits[0] == 0},
	}
}

func (proc *process) compute(wg *sync.WaitGroup, args []int) {
	defer close(proc.stdout)
	defer wg.Done()
	for {
		proc.decodeInstruction()
		switch proc.instr.op {
		case 1:
			proc.add()
		case 2:
			proc.mul()
		case 3:
			proc.read(&args)
		case 4:
			proc.write()
		case 5:
			proc.jnz()
		case 6:
			proc.jz()
		case 7:
			proc.lessThan()
		case 8:
			proc.equals()
		case 99:
			return
		default:
			log.Printf("Invalid opcode: %d\n", proc.instr.op)
			os.Exit(1)
		}
	}
}
