package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const A = 'A'
const B = 'B'
const C = 'C'
const L = 'L'
const R = 'R'
const SEP = ','
const END = '\n'

func (proc *process) calcAlignment() {
	sum := 0
	data := make([]string, 0)
	var line strings.Builder
	for v := range proc.stdout {
		if v == '\n' {
			data = append(data, line.String())
			line.Reset()
		} else {
			line.WriteByte(byte(v))
		}
	}
	for i := 1; i < len(data)-2; i++ {
		line := data[i]
		for j := 1; j < len(line)-2; j++ {
			if data[i][j] != '#' {
				continue
			}
			if data[i+1][j] == '#' && data[i-1][j] == '#' && data[i][j+1] == '#' && data[i][j-1] == '#' {
				sum += i * j
			}
		}
	}
	fmt.Println(strings.Join(data, "\n"))
	fmt.Println("Alignment parameters:", sum)
}

func (proc *process) dustingReport() {
	for v := range proc.stdout {
		if v > '~' {
			fmt.Printf("Dust collected: %d\n", v)
		}
	}
}

func runASCII(code intCode, mode int64) {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	proc.code[0] = mode
	go proc.compute([]int64{
		// Main movement routine
		A, SEP, A, SEP, B, SEP, C, SEP, C, SEP, A, SEP, C, SEP, B, SEP, C, SEP, B, END,
		// Movement function A
		L, SEP, '4', SEP, L, SEP, '4', SEP, L, SEP, '6', SEP, R, SEP, '1', '0', SEP, L, SEP, '6', END,
		// Movement function B
		L, SEP, '1', '2', SEP, L, SEP, '6', SEP, R, SEP, '1', '0', SEP, L, SEP, '6', END,
		// Movement function C
		R, SEP, '8', SEP, R, SEP, '1', '0', SEP, L, SEP, '6', END,
		// Continuous video feed
		'n', END,
	})
	if mode == 1 {
		proc.calcAlignment()
	} else {
		proc.dustingReport()
	}
}

func main() {
	data, err := ioutil.ReadFile("d17/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	runASCII(code, 1)
	runASCII(code, 2)
}
