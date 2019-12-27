package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func sendInput(proc *process, text string) {
	for i := range text {
		proc.stdin <- int64(text[i])
	}
}

// Pick up all possible items manually and set up at the exit--brute forcing the pressure sensor is instantaneous
func playTheGame(code intCode) {
	defer func() {
		recover()
		os.Exit(0)
	}()
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	go proc.compute([]int64{})
	go func() {
		for v := range proc.stdout {
			fmt.Printf("%c", v)
		}
	}()
	pickedUp := make(map[string]bool)
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if strings.HasPrefix(text, "take") {
			pickedUp[strings.SplitN(text[:len(text)-1], " ", 2)[1]] = true
		} else if strings.HasPrefix(text, "drop") {
			delete(pickedUp, strings.SplitN(text[:len(text)-1], " ", 2)[1])
		} else if text == "bruteforce\n" {
			var items []string
			for item := range pickedUp {
				items = append(items, item)
			}
			for {
				for _, item := range items {
					sendInput(proc, fmt.Sprintf("drop %s\n", item))
				}
				n := (rand.Int() % len(items)) + 1
				for i := 0; i < n; i++ {
					sendInput(proc, fmt.Sprintf("take %s\n", items[rand.Int()%len(items)]))
				}
				sendInput(proc, "west\n")
			}
		}
		sendInput(proc, text)
	}
}

func getCode(cmds string) []int64 {
	code := make([]int64, len(cmds))
	for i := range cmds {
		code[i] = int64(cmds[i])
	}
	return code
}

func runSolution(code intCode) *process {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	go proc.compute(getCode(strings.Join(
		[]string{
			"south",
			"take fixed point",
			"north",
			"north",
			"take candy cane",
			"west",
			"west",
			"take shell",
			"east",
			"east",
			"north",
			"north",
			"take polygon",
			"south",
			"west",
			"west",
			"west\n",
		}, "\n")))
	for v := range proc.stdout {
		fmt.Printf("%c", v)
	}
	return proc
}

func main() {
	data, err := ioutil.ReadFile("d25/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	//playTheGame(code)
	runSolution(code)
}
