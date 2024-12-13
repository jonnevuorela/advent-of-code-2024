package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type machine struct {
	A     [2]int
	B     [2]int
	prize [2]int
}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/13/input")
	machines := parseData(data)
	price := 0
	for i := range machines {
		price += calculatePrice(machines[i])

		fmt.Printf("%d machines out of %d processed", i+1, len(machines))
		fmt.Println()

	}
	fmt.Printf("Total price of tokens %d", price)

}
func calculatePrice(machine machine) int {

	// combinations of a and b up to 100 each
	for a := 0; a <= 100; a++ {
		for b := 0; b <= 100; b++ {
			posX := a*machine.A[0] + b*machine.B[0]
			posY := a*machine.A[1] + b*machine.B[1]

			if posX == machine.prize[0] && posY == machine.prize[1] {
				return (a * 3) + (b * 1)
			}
		}
	}

	return 0 // no solution
}

func splitByBias(value float64, bias float64) (float64, float64) {
	return value * bias, value * (1 - bias)
}

func parseData(data []byte) []machine {
	machines := []machine{}
	groups := bytes.Split(data, []byte("\n\n"))

	for i := range groups {
		lines := bytes.Split(bytes.TrimSpace(groups[i]), []byte("\n"))
		if len(lines) < 3 {
			continue
		}

		nm := machine{}

		// parse A
		aLine := string(lines[0])
		aParts := strings.Split(strings.TrimPrefix(aLine, "Button A: X+"), ", Y+")
		if len(aParts) == 2 {
			x, _ := strconv.Atoi(aParts[0])
			y, _ := strconv.Atoi(aParts[1])
			nm.A = [2]int{x, y}
		}

		// parse B
		bLine := string(lines[1])
		bParts := strings.Split(strings.TrimPrefix(bLine, "Button B: X+"), ", Y+")
		if len(bParts) == 2 {
			x, _ := strconv.Atoi(bParts[0])
			y, _ := strconv.Atoi(bParts[1])
			nm.B = [2]int{x, y}
		}

		// parse prize
		prizeLine := string(lines[2])
		prizeParts := strings.Split(strings.TrimPrefix(prizeLine, "Prize: X="), ", Y=")
		if len(prizeParts) == 2 {
			x, _ := strconv.Atoi(prizeParts[0])
			y, _ := strconv.Atoi(prizeParts[1])
			nm.prize = [2]int{x, y}
		}

		machines = append(machines, nm)
	}
	return machines
}
