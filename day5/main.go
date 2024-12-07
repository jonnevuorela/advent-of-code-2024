package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/5/input")

	input := bytes.Split(data, []byte("\n"))

	pattern := `\d{2}\|\d{2}`
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	ins := extractInstructions(input, *re)
	u := extractUpdates(input, *re)
	result := validateUpdate(u, ins)
	fmt.Println(result)
}

func validateUpdate(u [][]byte, ins [][]byte) int {
	count := 0
	validUpdates := [][]byte{}

	for _, update := range u {
		numbers := bytes.Split(update, []byte(","))
		for j := range numbers {
			numbers[j] = bytes.Trim(numbers[j], ",")
		}

		fmt.Printf("update: %s\n", update)
		if isValidUpdate(numbers, ins) {
			validUpdates = append(validUpdates, update)
			middleIndex := len(numbers) / 2
			middleNum, err := strconv.Atoi(string(numbers[middleIndex]))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Valid update: %s \n", update)
			count += middleNum
		} else {
			fmt.Printf("Invalid update: %s\n", update)
		}
	}

	fmt.Printf("Total valid updates: %d\n", len(validUpdates))
	return count
}

func isValidUpdate(numbers [][]byte, ins [][]byte) bool {
	for _, rule := range ins {
		parts := bytes.Split(rule, []byte("|"))
		firstIns, insErr1 := strconv.Atoi(string(parts[0]))
		secondIns, insErr2 := strconv.Atoi(string(parts[1]))
		if insErr1 != nil || insErr2 != nil {
			log.Fatal(insErr1, insErr2)
		}

		firstIndex := -1
		secondIndex := -1
		for idx, num := range numbers {
			numInt, err := strconv.Atoi(string(num))
			if err != nil {
				log.Fatal(err)
			}
			if numInt == firstIns {
				firstIndex = idx
			}
			if numInt == secondIns {
				secondIndex = idx
			}
		}

		if firstIndex != -1 && secondIndex != -1 {
			if firstIndex > secondIndex {
				fmt.Printf("Invalid due to rule: %s\n", rule)
				return false
			}
		}
	}
	return true
}

func extractUpdates(input [][]byte, re regexp.Regexp) [][]byte {
	var updates [][]byte
	for i := 0; i < len(input); i++ {
		if !re.MatchString(string(input[i])) && len(input[i]) != 0 {
			updates = append(updates, input[i])
		}
	}
	return updates
}

func extractInstructions(input [][]byte, re regexp.Regexp) [][]byte {
	var instructions [][]byte
	for i := 0; i < len(input); i++ {
		if re.MatchString(string(input[i])) {
			instructions = append(instructions, input[i])
		}
	}
	return instructions
}
