package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const length int = 6
const width int = 25

type (
	layer struct {
		data   string
		digits map[int32]int
	}

	image []layer
)

func (layer *layer) countDigits() {
	layer.digits = make(map[int32]int)
	for _, digit := range layer.data {
		layer.digits[digit]++
	}
}

func (image image) getChecksum() int {
	var minZeroes map[int32]int
	for i := range image {
		if minZeroes == nil || image[i].digits['0'] < minZeroes['0'] {
			minZeroes = image[i].digits
		}
	}
	return minZeroes['1'] * minZeroes['2']
}

func (image image) decode() {
	res := image[0].data
	for i := 0; i < len(image); i++ {
		data := image[i].data
		for j := 0; j < len(res); j++ {
			if res[j] == '2' {
				res = res[:j] + data[j:j+1] + res[j+1:]
			}
		}
	}
	for row := 0; row < length; row++ {
		fmt.Println(strings.ReplaceAll(res[:width], "0", " "))
		res = res[width:]
	}
}

func main() {
	data, err := ioutil.ReadFile("d08/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := string(data)
	size := length * width
	image := make(image, len(input)/size)
	for i := range image {
		image[i].data = input[:size]
		image[i].countDigits()
		input = input[size:]
	}
	fmt.Println("Checksum:", image.getChecksum())
	image.decode()
}
