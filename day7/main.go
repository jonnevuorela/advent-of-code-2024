package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"math"
	"strconv"
)

type operator int

const (
	Add operator = iota
	Multiply
	Concatenate
)

type eq struct {
	ans   []byte
	terms [][]byte
}

var count int

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/7/input")

	equations := extractEquation(data)
	for i := range equations {
		figureAnswer(equations[i])
	}
	fmt.Println()
	figureAnswer(equations[0])
	fmt.Printf("sum of answers %d", count)
}
func figureAnswer(eq eq) {
	numOperators := len(eq.terms) - 1
	numCombinations := int(math.Pow(3, float64(numOperators)))
	expectedAns, _ := strconv.Atoi(string(eq.ans))
	terms := convertTermsToInt(eq.terms)

	for combination := 0; combination < numCombinations; combination++ {
		operators := make([]operator, numOperators)

		for i := 0; i < numOperators; i++ {
			switch (combination / int(math.Pow(3, float64(i)))) % 3 {
			case 0:
				operators[i] = Add
			case 1:
				operators[i] = Multiply
			case 2:
				operators[i] = Concatenate
			}
		}

		result := terms[0]
		for i := 0; i < numOperators; i++ {
			switch operators[i] {
			case Add:
				result += terms[i+1]
			case Multiply:
				result *= terms[i+1]
			case Concatenate:
				result, _ = strconv.Atoi(fmt.Sprintf("%d%d", result, terms[i+1]))
			}
		}

		if result == expectedAns {
			count += result
			return
		}
	}
}

func convertTermsToInt(terms [][]byte) []int {
	intTerms := make([]int, len(terms))
	for i, term := range terms {
		intTerms[i], _ = strconv.Atoi(string(term))
	}
	return intTerms
}
func extractEquation(data []byte) []eq {
	input := bytes.Split(data, []byte("\n"))
	equations := []eq{}
	for i := range input {
		if len(input[i]) == 0 {
			continue
		}

		e := eq{}
		parts := bytes.Split(input[i], []byte(":"))
		e.ans = parts[0]
		e.terms = bytes.Split(parts[1], []byte(" "))
		equations = append(equations, e)
	}
	return equations
}
