package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"time"
)

type location struct {
	char        byte
	location    [2]int
	highlight   bool
	hasAntenna  bool
	hasAntinode bool
}

var count int
var finalCount int

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/8/input")
	md := mapInput(data)

	lastPrinted := make([][]byte, len(md))
	for i := range lastPrinted {
		lastPrinted[i] = make([]byte, len(md[i]))
		for j := range lastPrinted[i] {
			lastPrinted[i][j] = md[i][j].char
		}
	}

	// Initialize the map
	var buffer bytes.Buffer
	initializeMap(&buffer, md)
	fmt.Print(buffer.String())

	ticker := time.Tick(time.Second / 100)
	running := true
	currentRow, currentCol := 0, 0

	for running {
		for range ticker {
			if currentRow < len(md) && currentCol < len(md[currentRow]) {
				md[currentRow][currentCol].highlight = true
				go removeHighlight(&md[currentRow][currentCol])

				if md[currentRow][currentCol].hasAntenna {
					matches := checkMatches(md, md[currentRow][currentCol])
					for _, match := range matches {
						md[match.location[1]-1][match.location[0]-1].highlight = true
						go removeHighlight(&md[match.location[1]-1][match.location[0]-1])

						xDiff, yDiff := checkDistance(md[currentRow][currentCol], match)
						setAntinode(md, md[currentRow][currentCol], match, xDiff, yDiff)
					}
				}

				// Move to the next location
				currentCol++
				if currentCol >= len(md[currentRow]) {
					currentCol = 0
					currentRow++
					if currentRow >= len(md) {
						// Stop the ticker and remove all highlights before exiting
						running = false
						for i := range md {
							for j := range md[i] {
								md[i][j].highlight = false
								if md[i][j].hasAntinode {
									finalCount++
								}
							}
						}
						buffer.Reset()
						updateMap(&buffer, md, lastPrinted)
						fmt.Print(buffer.String())
						fmt.Printf("\033[%d;0H\033[K", len(md)+3)
						fmt.Printf("Antinodes found: %d\n", finalCount)
						break
					}
				}
				buffer.Reset()
				updateMap(&buffer, md, lastPrinted)
				fmt.Print(buffer.String())

				// Move the cursor to the position below the graphics and clear the line
				fmt.Printf("\033[%d;0H\033[K", len(md)+2)

				// Use the new function to count all antinodes
				fmt.Printf("Antinodes found so far: %d\n", count)
			}
		}
	}
}

func updateMap(buffer *bytes.Buffer, md [][]location, lastPrinted [][]byte) {
	for i := range md {
		rowChanged := false
		for j := range md[i] {
			if md[i][j].char != lastPrinted[i][j] || md[i][j].highlight {
				rowChanged = true
			}
		}

		// If the row has changed, update the entire row
		if rowChanged {
			buffer.WriteString(fmt.Sprintf("\033[%d;0H", i+1))
			for j := range md[i] {
				if md[i][j].highlight {
					buffer.WriteString(fmt.Sprintf("\033[38;2;255;255;102m%c \033[0m", md[i][j].char)) // Yellow color for highlighted
				} else if md[i][j].char == '#' {
					buffer.WriteString(fmt.Sprintf("\033[31m%c \033[0m", md[i][j].char)) // Red color for '#'
				} else {
					buffer.WriteString(fmt.Sprintf("%c ", md[i][j].char))
				}
				lastPrinted[i][j] = md[i][j].char
			}
		}
	}
}

func removeHighlight(cell *location) {
	time.Sleep(200 * time.Millisecond)
	cell.highlight = false
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
				char:        char,
				hasAntinode: false,
				highlight:   false,
				hasAntenna:  false,
			}
			l.location[0] = j + 1
			l.location[1] = i + 1
			if l.char != byte('.') {
				l.hasAntenna = true
			}
			rowLocations = append(rowLocations, l)
		}
		locations = append(locations, rowLocations)
	}
	return locations
}

func checkMatches(md [][]location, loc location) []location {
	var matches []location
	for i := range md {
		for j := range md[i] {
			if md[i][j].char == loc.char && md[i][j].char != '#' && md[i][j].location != loc.location {
				matches = append(matches, md[i][j])
			}
		}
	}
	return matches
}

func checkDistance(loc1, loc2 location) (int, int) {
	xDiff := loc2.location[0] - loc1.location[0]
	yDiff := loc2.location[1] - loc1.location[1]
	return xDiff, yDiff
}

func setAntinode(md [][]location, loc1, loc2 location, xDiff, yDiff int) {
	inverseX1 := loc1.location[0] - xDiff
	inverseY1 := loc1.location[1] - yDiff

	// Adjust for zero-based indexing
	inverseX1 -= 1
	inverseY1 -= 1

	if inverseX1 >= 0 && inverseX1 < len(md[0]) && inverseY1 >= 0 && inverseY1 < len(md) {
		if xDiff != 0 || yDiff != 0 {
			if !md[inverseY1][inverseX1].hasAntinode {
				if !md[inverseY1][inverseX1].hasAntenna {
					md[inverseY1][inverseX1].char = '#'
					md[inverseY1][inverseX1].hasAntinode = true
					count++
				} else if md[inverseY1][inverseX1].hasAntenna {
					md[inverseY1][inverseX1].hasAntinode = true
					count++
				}
			}
		}
	}

	inverseX2 := loc2.location[0] + xDiff - 1
	inverseY2 := loc2.location[1] + yDiff - 1

	if inverseX2 >= 0 && inverseX2 < len(md[0]) && inverseY2 >= 0 && inverseY2 < len(md) {
		if xDiff != 0 || yDiff != 0 {
			if !md[inverseY2][inverseX2].hasAntinode {
				if !md[inverseY2][inverseX2].hasAntenna {
					md[inverseY2][inverseX2].char = '#'
					md[inverseY2][inverseX2].hasAntinode = true
					count++
				} else if md[inverseY2][inverseX2].hasAntenna {
					md[inverseY2][inverseX2].hasAntinode = true
					count++
				}
			}
		}
	}
}
