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
	data := input.GetInput("https://adventofcode.com/2024/day/3/input")
	pattern := regexp.MustCompile(`mul\(\d+,\d+\)`)
	parsedData := parseInput(data, *pattern)
	output := extractInput(parsedData)

	fmt.Println(calculateInput(output))

}
func calculateInput(input [][]int) int {
	var total int
	for i := 0; i < len(input); i++ {
		sum := input[i][0] * input[i][1]
		total += sum
	}
	return total

}
func extractInput(data [][]byte) [][]int {

	var output [][]int

	for i := 0; i < len(data); i++ {
		pattern := regexp.MustCompile(`\((\d+),(\d+)\)`)
		matches := pattern.FindSubmatch(data[i])

		if len(matches) == 3 {
			num1, err1 := strconv.Atoi(string(matches[1]))
			num2, err2 := strconv.Atoi(string(matches[2]))
			if err1 == nil && err2 == nil {
				pair := []int{num1, num2}
				output = append(output, pair)
			} else {
				log.Fatal(err1, err2)
			}
		} else {
			fmt.Println("Pattern not matched")
		}
	}
	return output

}

func parseInput(data []byte, pattern regexp.Regexp) [][]byte {
	var results []byte

	done := false

	for !done {
		loc := pattern.FindIndex(data)

		if loc != nil {
			result := (data[loc[0]:loc[1]])
			results = append(results, result...)

			data = append(data[:loc[0]], data[loc[1]:]...)
		} else {
			done = true
		}
	}
	resultss := bytes.Split(results, []byte("mul"))
	return resultss
}
