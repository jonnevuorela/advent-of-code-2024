package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"time"
)

type location struct {
	char       byte
	location   [2]int
	obstructed bool
	visited    bool
}
type direction struct {
	f [2]int
	r [2]int
	b [2]int
	l [2]int
}
type guard struct {
	location  [2]int
	direction direction
}

var g guard
var md [][]location
var running bool
var count int

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/6/input")
	g = newGuard()
	md = mapInput(data)

	ticker := time.Tick(time.Second / 100)
	lastPrinted := make([][]byte, len(md))
	for i := range lastPrinted {
		lastPrinted[i] = make([]byte, len(md[i]))
	}

	// Initialize the map
	var buffer bytes.Buffer
	initializeMap(&buffer, md)
	fmt.Print(buffer.String())

	running = true

	for running {
		for range ticker {
			buffer.Reset()
			updateMap(&buffer, md, lastPrinted, g)
			fmt.Print(buffer.String())

			fmt.Printf("\033[%d;0H", len(md)+1)
			fmt.Printf("guard location x: %d y %d", g.location[0], g.location[1])
			fmt.Printf("visited locations %d", count)

			simulateRoute(md)
		}
	}
}

func simulateRoute(md [][]location) {
	nextLocation := [2]int{
		g.location[0] + g.direction.f[0],
		g.location[1] + g.direction.f[1],
	}
	if nextLocation[0] < 0 || nextLocation[1] < 0 || nextLocation[0] > len(md)+1 || nextLocation[1] > len(md[0])+1 {
		count = calculateXBytes(md)
		return
	}
	var nextLoc *location
	for i := range md {
		for j := range md[i] {
			if md[i][j].location[0] == nextLocation[0] && md[i][j].location[1] == nextLocation[1] {
				nextLoc = &md[i][j]
				break
			}
		}
		if nextLoc != nil {
			break
		}
	}
	if nextLoc != nil {
		if nextLoc.obstructed {
			turnGuard()
		} else {
			g.location = nextLocation
		}
	} else {
		count = calculateXBytes(md)
	}
}

func turnGuard() {
	// Rotate the guard's direction clockwise
	g.direction.f, g.direction.r, g.direction.b, g.direction.l =
		g.direction.r, g.direction.b, g.direction.l, g.direction.f
}

func initializeMap(buffer *bytes.Buffer, md [][]location) {
	for i := range md {
		for j := range md[i] {
			// Move the cursor to the correct position
			buffer.WriteString(fmt.Sprintf("\033[%d;%dH", i+1, j*2+1))
			buffer.WriteString(fmt.Sprintf("%c ", md[i][j].char))
		}
	}
}

func updateMap(buffer *bytes.Buffer, md [][]location, lastPrinted [][]byte, guard guard) {
	for i := range md {
		rowChanged := false
		for j := range md[i] {
			// Check if the current location is the guard's location
			if md[i][j].location[0] == guard.location[0] && md[i][j].location[1] == guard.location[1] {
				md[i][j].char = '^'
				md[i][j].visited = true
			} else if !md[i][j].obstructed && md[i][j].visited {
				md[i][j].char = 'X'
			}

			// Check if the current character is different from the last printed character
			if lastPrinted[i][j] != md[i][j].char {
				rowChanged = true
			}
		}

		// If the row has changed, update the entire row
		if rowChanged {
			buffer.WriteString(fmt.Sprintf("\033[%d;0H", i+1))
			for j := range md[i] {
				if md[i][j].char == 'X' {
					buffer.WriteString(fmt.Sprintf("\033[31m%c \033[0m", md[i][j].char)) // Red color for 'X'
				} else {
					buffer.WriteString(fmt.Sprintf("%c ", md[i][j].char))
				}
				lastPrinted[i][j] = md[i][j].char
			}
		}
	}
}

func newGuard() guard {
	return guard{
		direction: direction{
			f: [2]int{0, -1},
			r: [2]int{1, 0},
			b: [2]int{0, 1},
			l: [2]int{-1, 0},
		},
	}
}

func mapInput(data []byte) [][]location {
	var locations [][]location
	rows := bytes.Split(data, []byte("\n"))
	for i, row := range rows {
		if len(row) == 0 {
			continue
		}
		var rowLocations []location
		for j, char := range row {
			l := location{
				char: char,
			}
			l.location[0] = j + 1
			l.location[1] = i + 1
			if l.char == byte('#') {
				l.obstructed = true
			} else if l.char == byte('^') {
				g.location[0] = l.location[0]
				g.location[1] = l.location[1]
				l.visited = true
				l.char = byte('.')
			} else {
				l.visited = false
			}
			rowLocations = append(rowLocations, l)
		}
		locations = append(locations, rowLocations)
	}
	return locations
}

func calculateXBytes(md [][]location) int {
	count := 0
	for i := range md {
		for j := range md[i] {
			if md[i][j].char == 'X' {
				count++
			}
		}
	}
	running = false
	return count
}
