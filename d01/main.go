package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func calculateFuel(moduleMasses []float64, countFuel bool) float64 {
	var total float64
	for _, mass := range moduleMasses {
		fuel := math.Floor(mass/3.0) - 2
		if countFuel {
			fuelMass := fuel
			for {
				extra := math.Floor(fuelMass/3.0) - 2
				if extra < 1.0 {
					break
				}
				fuel += extra
				fuelMass = extra
			}
		}
		total += fuel
	}
	return total
}

func main() {
	file, err := os.Open("d01/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var moduleMasses []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mass, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			log.Println(err)
			continue
		}
		moduleMasses = append(moduleMasses, mass)
	}
	fmt.Printf("Total fuel: %d\n", int64(calculateFuel(moduleMasses, false)))
	fmt.Printf("With fuel for fuel: %d\n", int64(calculateFuel(moduleMasses, true)))
}
