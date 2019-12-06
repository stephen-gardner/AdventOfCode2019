package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type (
	orbitals      map[string]string
	transferRoute map[string]bool
)

var orbitMap = make(orbitals)

func (orbitals orbitals) countTotalOrbits() int {
	total := 0
	for origin := range orbitals {
		for origin = orbitals[origin]; origin != ""; origin = orbitals[origin] {
			total++
		}
	}
	return total
}

func (orbitals orbitals) getRouteToEnd(satellite string) transferRoute {
	route := make(transferRoute)
	for origin := orbitals[satellite]; origin != ""; origin = orbitals[origin] {
		route[origin] = true
	}
	return route
}

func (orbitals orbitals) countTransfersToIntersection(route transferRoute, satellite string) int {
	transfers := 0
	for origin := orbitals[satellite]; ; origin = orbitals[origin] {
		if _, present := route[origin]; present {
			break
		}
		transfers++
	}
	return transfers
}

func (orbitals orbitals) countTransfersToSAN() int {
	return orbitals.countTransfersToIntersection(orbitals.getRouteToEnd("SAN"), "YOU") +
		orbitals.countTransfersToIntersection(orbitals.getRouteToEnd("YOU"), "SAN")
}

func main() {
	data, err := ioutil.ReadFile("d06/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	for _, raw := range lines {
		pair := strings.Split(raw, ")")
		orbitMap[pair[1]] = pair[0]
		if _, present := orbitMap[pair[0]]; !present {
			orbitMap[pair[0]] = ""
		}
	}
	fmt.Printf("Total Orbits: %d\n", orbitMap.countTotalOrbits())
	fmt.Printf("Transfers to SAN: %d\n", orbitMap.countTransfersToSAN())
}
