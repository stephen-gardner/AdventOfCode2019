package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var recipes = make(map[string]map[string]int)
var productAmounts = make(map[string]int)
var costs = make(map[string]float64)

func calculateOre(needed map[string]int) int {
	added := make(map[string]int)
	for more := true; more; {
		more = false
		for product, quantity := range needed {
			if product == "ORE" || needed[product] <= 0 {
				continue
			}
			more = true
			for ; quantity > 0; quantity -= productAmounts[product] {

				for reagent, amount := range recipes[product] {
					added[reagent] += amount
				}
				needed[product] -= productAmounts[product]
			}
		}
		for product, quantity := range added {
			needed[product] += quantity
			delete(added, product)
		}
	}
	return needed["ORE"]
}

func calculateEfficientCost(product string) float64 {
	if _, ok := costs[product]; ok {
		return costs[product]
	}
	ore := float64(0)
	recipe := recipes[product]
	for reagent, amount := range recipe {
		if reagent == "ORE" {
			ore += float64(amount)
			continue
		}
		ore += calculateEfficientCost(reagent) * float64(amount)
	}
	costs[product] = ore / float64(productAmounts[product])
	return costs[product]
}

func main() {
	data, err := ioutil.ReadFile("d14/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	for _, fline := range lines {
		formula := strings.Split(fline, "=>")
		productLine := strings.Trim(formula[1], " ")
		productPair := strings.Split(productLine, " ")
		product := productPair[1]
		amount, _ := strconv.Atoi(productPair[0])
		productAmounts[product] = amount
		ingredientsList := strings.Split(strings.Trim(formula[0], " "), ",")
		ingredientMap := make(map[string]int)
		for _, ingredient := range ingredientsList {
			pair := strings.Split(strings.Trim(ingredient, " "), " ")
			amount, _ := strconv.Atoi(pair[0])
			ingredientMap[pair[1]] = amount
		}
		recipes[product] = ingredientMap
	}
	oreCost := calculateOre(map[string]int{"FUEL": 1})
	fmt.Println("Ore needed:", oreCost)
	// This probably only works because the problem was designed to consume
	// just enough materials to produce enough fuel to consume a trillion ore
	fmt.Println("Fuel for 1 trillion ore:", int(1000000000000.0/calculateEfficientCost("FUEL")))
}
