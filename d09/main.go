package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func runBOOST(code intCode, mode int64) {
	proc := &process{
		code:   append(code.copy(), make(intCode, 4096)...),
		stdin:  make(chan int64),
		stdout: make(chan int64),
	}
	go proc.compute([]int64{mode})
	for out := range proc.stdout {
		fmt.Println(out)
	}
}

func main() {
	data, err := ioutil.ReadFile("d09/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	strValues := strings.Split(string(data), ",")
	code := make(intCode, len(strValues))
	for i := range strValues {
		code[i], _ = strconv.ParseInt(strValues[i], 10, 64)
	}
	fmt.Printf("BOOST test keycode: ")
	runBOOST(code, 1)
	fmt.Printf("Distress signal coordinates: ")
	runBOOST(code, 2)
}
