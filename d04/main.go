package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type password []int

func getPassword(n int) password {
	numDigits := 1
	for tmp := n; tmp >= 10; numDigits++ {
		tmp /= 10
	}
	res := make([]int, numDigits)
	for i := numDigits - 1; i >= 0; i-- {
		res[i] = n % 10
		n /= 10
	}
	return res
}

func (pw password) getRunLength(i int) int {
	run := 0
	digit := pw[i]
	for ; i < len(pw); i++ {
		if pw[i] == digit {
			run++
		} else {
			break
		}
	}
	return run
}

func (pw password) hasDouble(exclusivePair bool) bool {
	for i := 1; i < len(pw); i++ {
		if pw[i] == pw[i-1] {
			if exclusivePair {
				runLength := pw.getRunLength(i - 1)
				if runLength == 2 {
					return true
				}
				i += runLength - 1
			} else {
				return true
			}
		}
	}
	return false
}

func (pw password) increases() bool {
	for i := 1; i < len(pw); i++ {
		if pw[i] < pw[i-1] {
			return false
		}
	}
	return true
}

func main() {
	data, err := ioutil.ReadFile("d04/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	rangeStr := strings.Split(string(data), "-")
	passwordRange := make([]int, 2)
	passwordRange[0], _ = strconv.Atoi(rangeStr[0])
	passwordRange[1], _ = strconv.Atoi(rangeStr[1])
	allowed := 0
	allowedExclusive := 0
	for n := passwordRange[0]; n <= passwordRange[1]; n++ {
		pw := getPassword(n)
		if pw.hasDouble(false) && pw.increases() {
			allowed++
			if pw.hasDouble(true) {
				allowedExclusive++
			}
		}
	}
	fmt.Printf("Acceptable permutations: %d\nWith exclusive pair: %d\n", allowed, allowedExclusive)
}
