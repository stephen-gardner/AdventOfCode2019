package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func newStack(cards []int) []int {
	res := make([]int, len(cards))
	for i := 0; i < len(cards); i++ {
		res[i] = cards[len(cards)-i-1]
	}
	return res
}

func cut(cards []int, n int) []int {
	if n > 0 {
		return append(cards[n:], cards[:n]...)
	} else {
		n *= -1
		return append(cards[len(cards)-n:], cards[:len(cards)-n]...)
	}
}

func increment(cards []int, n int) []int {
	res := make([]int, len(cards))
	i := 0
	for _, v := range cards {
		res[i] = v
		i = (i + n) % len(cards)
	}
	return res
}

func findCard(cards []int, v int) int {
	for i := range cards {
		if cards[i] == v {
			return i
		}
	}
	return -1
}

func process(cards []int, instr []string) []int {
	for _, op := range instr {
		if op == "deal into new stack" {
			cards = newStack(cards)
		} else if strings.HasPrefix(op, "cut") {
			words := strings.Split(op, " ")
			n, _ := strconv.Atoi(words[1])
			cards = cut(cards, n)
		} else if strings.HasPrefix(op, "deal with increment") {
			words := strings.Split(op, " ")
			n, _ := strconv.Atoi(words[3])
			cards = increment(cards, n)
		}
	}
	return cards
}

func main() {
	input, err := ioutil.ReadFile("d22/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	instr := strings.Split(string(input), "\n")
	cards := make([]int, 10007)
	for i := range cards {
		cards[i] = i
	}
	cards = process(cards, instr)
	fmt.Println("Position of card 2019:", findCard(cards, 2019))
}
