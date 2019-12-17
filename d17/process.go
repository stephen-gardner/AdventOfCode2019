package main

import (
	"fmt"
	"log"
)

type (
	intCode []int64

	process struct {
		code   intCode
		stdin  chan int64
		stdout chan int64
		RBP    int64
		PC     int64
	}
)

func (code intCode) copy() intCode {
	return append(intCode{}, code...)
}

func (code intCode) dump() {
	fmt.Printf("%v\n", code)
}

// Opcode #1
func (proc *process) add() {
	v := proc.getOperands(false, false, true)
	proc.code[v[2]] = v[0] + v[1]
	proc.PC += 4
}

// Opcode #2
func (proc *process) mul() {
	v := proc.getOperands(false, false, true)
	proc.code[v[2]] = v[0] * v[1]
	proc.PC += 4
}

// Opcode #3
func (proc *process) read(args *[]int64) {
	var val int64
	if len(*args) > 0 {
		val = (*args)[0]
		*args = (*args)[1:]
	} else {
		val = <-proc.stdin
	}
	proc.code[proc.getOperands(true)[0]] = val
	proc.PC += 2
}

// Opcode #4
func (proc *process) write() {
	proc.stdout <- proc.getOperands(false)[0]
	proc.PC += 2
}

// Opcode #5
func (proc *process) jnz() {
	v := proc.getOperands(false, false)
	if v[0] != 0 {
		proc.PC = v[1]
		return
	}
	proc.PC += 3
}

// Opcode #6
func (proc *process) jz() {
	v := proc.getOperands(false, false)
	if v[0] == 0 {
		proc.PC = v[1]
		return
	}
	proc.PC += 3
}

// Opcode #7
func (proc *process) lessThan() {
	v := proc.getOperands(false, false, true)
	if v[0] < v[1] {
		proc.code[v[2]] = 1
	} else {
		proc.code[v[2]] = 0
	}
	proc.PC += 4
}

// Opcode #8
func (proc *process) equals() {
	v := proc.getOperands(false, false, true)
	if v[0] == v[1] {
		proc.code[v[2]] = 1
	} else {
		proc.code[v[2]] = 0
	}
	proc.PC += 4
}

// Opcode #9
func (proc *process) offsetRBP() {
	proc.RBP += proc.getOperands(false)[0]
	proc.PC += 2
}

func (proc *process) decodeInstruction(nOperands int) []int64 {
	n := proc.code[proc.PC]
	digits := make([]int64, nOperands+2)
	for i := 0; i < nOperands+2; i++ {
		digits[i] = n % 10
		n /= 10
	}
	return digits[2:]
}

func (proc *process) getOperands(isAddress ...bool) []int64 {
	vars := make([]int64, len(isAddress))
	modes := proc.decodeInstruction(len(isAddress))
	for n := int64(0); n < int64(len(isAddress)); n++ {
		vars[n] = proc.code[proc.PC+n+1]
		switch modes[n] {
		case 0:
		case 1:
			continue
		case 2:
			vars[n] = proc.RBP + vars[n]
		default:
			log.Fatalf("Invalid operand mode: %d\n", modes[n])
		}
		if !isAddress[n] {
			vars[n] = proc.code[vars[n]]
		}
	}
	return vars
}

func (proc *process) compute(args []int64) {
	defer close(proc.stdout)
	for {
		switch op := proc.code[proc.PC] % 100; op {
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
		case 9:
			proc.offsetRBP()
		case 99:
			return
		default:
			log.Fatalf("Invalid opcode: %d\n", op)
		}
	}
}
