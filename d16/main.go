package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
)

var base = []int{0, 1, 0, -1}

func calcEntireSignal(data []int8, nPhases int) []int8 {
	res := make([]int8, len(data))
	for phase := 0; phase < nPhases; phase++ {
		for round := 0; round < len(data); round++ {
			digit := 0
			pos := 0
			if round == 0 {
				pos = 1
			}
			count := 1
			for _, n := range data {
				digit += int(n) * base[pos]
				if count >= round {
					pos = (pos + 1) % len(base)
					count = 0
				} else {
					count++
				}
			}
			digit = int(math.Abs(float64(digit))) % 10
			res[round] = int8(digit)
		}
		data = res
		res = make([]int8, len(data))
	}
	return data[:8]
}

func calcRelevantSignal(data []int8, nPhases int) []int8 {
	res := make([]int8, len(data))
	for phase := 0; phase < nPhases; phase++ {
		sum := 0
		for round := len(data) - 1; round >= 0; round-- {
			sum += int(data[round])
			res[round] = int8(sum % 10)
		}
		data = res
		res = make([]int8, len(data))
	}
	return data[:8]
}

func getOffset(data []int8) int {
	offset := 0
	for i := 0; i < 7; i++ {
		offset = (offset * 10) + int(data[i])
	}
	return offset
}

func main() {
	input, err := ioutil.ReadFile("d16/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	var data []int8
	for _, d := range string(input) {
		val, _ := strconv.Atoi(string(d))
		data = append(data, int8(val))
	}
	realData := make([]int8, len(data)*10000)
	for i := 0; i < len(realData); i += len(data) {
		for j := 0; j < len(data); j++ {
			realData[i+j] = data[j]
		}
	}
	realData = realData[getOffset(data):]
	fmt.Println("Test Signal:", calcEntireSignal(data, 100))
	fmt.Println("Real Signal:", calcRelevantSignal(realData, 100))
}
