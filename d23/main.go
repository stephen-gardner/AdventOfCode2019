package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const NAT = 255

func isIdle(idleStatuses map[int64]bool) bool {
	for _, idle := range idleStatuses {
		if !idle {
			return false
		}
	}
	return true
}

func route(computers []*process) {
	queue := make(map[int64][]int64)
	idle := make(map[int64]bool)
	seen := make(map[int64]bool)
	for i := int64(0); i < 50; i++ {
		queue[i] = []int64{}
	}
	for {
		for i := int64(0); i < 50; i++ {
			comp := computers[i]
			if i == 0 && isIdle(idle) && len(queue[NAT]) == 2 {
				y := queue[NAT][1]
				if _, dupe := seen[y]; dupe {
					fmt.Printf("NAT: %v*\n", queue[NAT])
					os.Exit(0)
				} else {
					fmt.Println("NAT:", queue[NAT])
				}
				seen[y] = true
				queue[0] = append(queue[0], queue[NAT]...)
				queue[NAT] = []int64{}
			}
			select {
			case address := <-comp.stdout:
				idle[i] = false
				data := []int64{<-comp.stdout, <-comp.stdout}
				if address == NAT {
					queue[NAT] = data
				} else {
					queue[address] = append(queue[address], data...)
				}
			default:
				v := int64(-1)
				hasIncoming := false
				if len(queue[i]) > 0 {
					v = queue[i][0]
					hasIncoming = true
				}
				select {
				case comp.stdin <- v:
					if hasIncoming {
						queue[i] = queue[i][1:]
					} else {
						idle[i] = true
					}
				default:
				}
			}
		}
	}
}

func startComputer(code intCode, address int64) *process {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	go proc.compute([]int64{address})
	return proc
}

func main() {
	data, err := ioutil.ReadFile("d23/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	computers := make([]*process, 50)
	for i := int64(0); i < 50; i++ {
		computers[i] = startComputer(code, i)
	}
	route(computers)
}
