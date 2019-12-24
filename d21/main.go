package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var commands = []string{
	"AND A J", "AND A T",
	"AND B J", "AND B T",
	"AND C J", "AND C T",
	"AND D J", "AND D T",
	"AND T J", "AND T T",
	"AND J J", "AND J T",
	"OR A J", "OR A T",
	"OR B J", "OR B T",
	"OR C J", "OR C T",
	"OR D J", "OR D T",
	"OR T J", "OR T T",
	"OR J J", "OR J T",
	"NOT A J", "NOT A T",
	"NOT B J", "NOT B T",
	"NOT C J", "NOT C T",
	"NOT D J", "NOT D T",
	"NOT T J", "NOT T T",
	"NOT J J", "NOT J T",
}

var extendedCommands = []string{
	"AND E J", "AND E T",
	"AND F J", "AND F T",
	"AND G J", "AND G T",
	"AND H J", "AND H T",
	"AND I J", "AND I T",
	"OR E J", "OR E T",
	"OR F J", "OR F T",
	"OR G J", "OR G T",
	"OR H J", "OR H T",
	"OR I J", "OR I T",
	"NOT E J", "NOT E T",
	"NOT F J", "NOT F T",
	"NOT G J", "NOT G T",
	"NOT H J", "NOT H T",
	"NOT I J", "NOT I T",
}

func stringToCode(str string) []int64 {
	res := make([]int64, len(str)+1)
	for i := range str {
		res[i] = int64(str[i])
	}
	res[len(str)] = '\n'
	return res
}

func getCode(instrStr []string) []int64 {
	var code []int64
	for _, str := range instrStr {
		code = append(code, stringToCode(str)...)
	}
	return code
}

func bruteForce(code intCode, run bool) {
	var wg sync.WaitGroup
	count := 0
	cmds := commands
	if run {
		cmds = append(cmds, extendedCommands...)
	}
	for nInstr := 1; nInstr < 15; nInstr++ {
		fmt.Println("Number of instructions:", nInstr)
		instrIdx := make([]int, nInstr+1)
		instrIdx[nInstr-1] = -1
		for {
			exit := false
			for i := nInstr - 1; i >= 0; i-- {
				instrIdx[i]++
				if instrIdx[i] >= len(cmds) {
					j := i
					for instrIdx[j] >= len(cmds) {
						if j == 0 {
							exit = true
							break
						}
						instrIdx[j] = 0
						instrIdx[j-1]++
						j--
					}
				} else {
					break
				}
			}
			if exit {
				break
			}
			instrStr := make([]string, len(instrIdx))
			for i := 0; i < nInstr; i++ {
				instrStr[i] = cmds[instrIdx[i]]
			}
			if run {
				instrStr[nInstr] = "RUN"
			} else {
				instrStr[nInstr] = "WALK"
			}
			count++
			wg.Add(1)
			go func(instr []int64) {
				proc := &process{
					code:   append(code.copy(), make(intCode, 4096)...),
					stdin:  make(chan int64),
					stdout: make(chan int64),
				}
				go proc.compute(instr)
				for v := range proc.stdout {
					if v > '~' {
						fmt.Printf("Solution found:\n%s\n", strings.Join(instrStr, "\n"))
						fmt.Println("Hull damage:", v)
						os.Exit(0)
					}
				}
				wg.Done()
			}(getCode(instrStr))
			if count == 100 {
				wg.Wait()
				count = 0
			}
		}
	}
}

func runSpringdroid(code intCode, instrStr []string) {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	go proc.compute(getCode(instrStr))
	for v := range proc.stdout {
		if v > '~' {
			if instrStr[len(instrStr)-1] == "RUN" {
				fmt.Println("Hull damage (extended range):", v)
			} else {
				fmt.Println("Hull damage:", v)
			}
		}
	}
}

func main() {
	data, err := ioutil.ReadFile("d21/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	runSpringdroid(code, []string{
		"OR A J",
		"AND C J",
		"NOT J J",
		"AND D J",
		"WALK",
	})
	runSpringdroid(code, []string{
		"NOT B J",
		"NOT C T",
		"OR T J",
		"AND H J",
		"NOT A T",
		"OR T J",
		"AND D J",
		"RUN",
	})
}
