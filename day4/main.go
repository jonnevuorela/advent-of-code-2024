package main

import (
	"aoc/input"
	"bytes"
	"fmt"
)

type letter struct {
	letter byte
	x      int
	y      int
}

func main() {
	byte := input.GetInput("https://adventofcode.com/2024/day/4/input")
	letters := mapInput(byte)
	count := handleCrossMas(letters)
	fmt.Printf("word count is: %d", count)
}

func mapInput(data []byte) [][]letter {
	var letters [][]letter
	rows := bytes.Split(data, []byte("\n"))
	for i, row := range rows {
		if len(row) == 0 {
			continue
		}
		var rowLetters []letter
		for j, char := range row {
			l := letter{
				letter: char,
				x:      j + 1,
				y:      i + 1,
			}
			rowLetters = append(rowLetters, l)
		}
		letters = append(letters, rowLetters)
	}
	return letters
}

func handleCrossMas(letters [][]letter) int {
	count := 0
	for i := range letters {
		for j := range letters[i] {
			if letters[i][j].letter == byte('A') {
				foundPattern := false
				// Check in the positive direction
				for k := 1; k < len(letters); k++ {
				}
				if foundPattern {
					count++
				}
			}
		}
	}
	return count
}

func handleXmas(letters [][]letter) int {
	count := 0
	word := []byte("XMAS")
	directions := [][2]int{
		{0, 1},   // horizontal forward
		{0, -1},  // horizontal backward
		{1, 0},   // vertical down
		{-1, 0},  // vertical up
		{1, 1},   // diagonal forward down
		{-1, -1}, // diagonal backward up
		{1, -1},  // diagonal forward up
		{-1, 1},  // diagonal backward down
	}

	for i := range letters {
		for j := range letters[i] {
			if letters[i][j].letter == byte('X') {
				for _, dir := range directions {
					if checkForWord(letters[i][j], word, letters, dir) {
						count++
					}
				}
			}
		}
	}
	return count
}

func checkForWord(letter letter, word []byte, letters [][]letter, dir [2]int) bool {
	wordLen := len(word)
	gridWidth := len(letters[0])
	gridHeight := len(letters)

	for i := 0; i < wordLen; i++ {
		newX := letter.x + dir[0]*i
		newY := letter.y + dir[1]*i

		if newX < 1 || newX > gridWidth || newY < 1 || newY > gridHeight {
			return false
		}

		if letters[newY-1][newX-1].letter != word[i] {
			return false
		}
	}
	return true
}
