package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type cell struct {
	char     byte
	location [2]int
	quadrant int
}
type robot struct {
	pos [2]int
	mag [2]int
}

func main() {
	fmt.Print("\033[0m")
	data := input.GetInput("https://adventofcode.com/2024/day/14/input")
	robots := mapRobots(data)
	area := mapArea()

	ticker := time.Tick(time.Second / 1000)
	lastPrinted := make([][]byte, len(area))
	for i := range lastPrinted {
		lastPrinted[i] = make([]byte, len(area[i]))
	}

	var buffer bytes.Buffer
	fmt.Print(buffer.String())

	sec := 0
	for range ticker {
		buffer.Reset()
		updateMap(&buffer, area, lastPrinted, robots)
		fmt.Print(buffer.String())
		for i := range robots {
			simulateRoute(&robots[i])
		}
		fmt.Printf("\033[%d;0H", len(area)+1)

		sec++
		if sec >= 100 {
			break
		}
	}
	calculateRobots(area, robots)

}

func calculateRobots(area [][]cell, robots []robot) {
	q1, q2, q3, q4 := 0, 0, 0, 0

	for _, robot := range robots {
		for i := range area {
			for j := range area[i] {
				if area[i][j].location[0] == robot.pos[0] &&
					area[i][j].location[1] == robot.pos[1] {
					switch area[i][j].quadrant {
					case 1:
						q1++
					case 2:
						q2++
					case 3:
						q3++
					case 4:
						q4++
					}
					break
				}
			}
		}
	}

	fmt.Println()

	count := q1 * q2 * q3 * q4
	fmt.Printf("Final Answer: %d", count)
}

func simulateRoute(r *robot) {
	nextLocation := [2]int{
		r.pos[0] + r.mag[0],
		r.pos[1] + r.mag[1],
	}

	if nextLocation[0] < 0 {
		nextLocation[0] = 101
	} else if nextLocation[0] > 101 {
		nextLocation[0] = 0
	}

	if nextLocation[1] < 0 {
		nextLocation[1] = 103
	} else if nextLocation[1] > 103 {
		nextLocation[1] = 0
	}

	r.pos = nextLocation
}

func updateMap(buffer *bytes.Buffer, md [][]cell, lastPrinted [][]byte, robots []robot) {
	for i := range md {
		for j := range md[i] {
			if i != 52 && j != 51 {
				if md[i][j].char == 'ø' {
					md[i][j].char = '.'
				}
			} else {

				if md[i][j].char == 'ø' {
					md[i][j].char = ' '
				}
			}
		}
	}

	for _, robot := range robots {
		for i := range md {
			for j := range md[i] {
				if md[i][j].location[0] == robot.pos[0] && md[i][j].location[1] == robot.pos[1] {
					md[i][j].char = 'ø'
				}
			}
		}
	}

	for i := range md {
		rowChanged := false
		for j := range md[i] {
			if lastPrinted[i][j] != md[i][j].char {
				rowChanged = true
				break
			}
		}

		if rowChanged {
			buffer.WriteString(fmt.Sprintf("\033[%d;0H", i+1))
			for j := range md[i] {
				if md[i][j].char == 'ø' {
					buffer.WriteString(fmt.Sprintf("\033[31m%c \033[0m", md[i][j].char))
				} else {
					buffer.WriteString(fmt.Sprintf("%c ", md[i][j].char))
				}
				lastPrinted[i][j] = md[i][j].char
			}
		}
	}
}

func mapRobots(data []byte) []robot {
	robots := []robot{}
	elems := bytes.Split(data, []byte("\n"))
	for i := range elems {
		if len(elems[i]) == 0 {
			continue
		}
		r := robot{}
		parts := bytes.Split(elems[i], []byte(" "))
		posParts := bytes.Split(parts[0], []byte(","))
		magParts := bytes.Split(parts[1], []byte(","))
		r.pos[0], _ = strconv.Atoi(string(bytes.TrimPrefix(posParts[0], []byte("p="))))
		r.pos[1], _ = strconv.Atoi(string(posParts[1]))
		r.mag[0], _ = strconv.Atoi(string(bytes.TrimPrefix(magParts[0], []byte("v="))))
		r.mag[1], _ = strconv.Atoi(string(magParts[1]))

		robots = append(robots, r)
	}
	return robots
}

func mapArea() [][]cell {
	area := [][]cell{}
	for i := 0; i < 103; i++ {
		row := []cell{}
		for j := 0; j < 101; j++ {
			cell := cell{}
			if i != 52 && j != 51 {
				cell.char = '.'
			} else {
				cell.char = ' '
			}
			cell.location[0] = j
			cell.location[1] = i
			if i < 52 && j < 51 {
				cell.quadrant = 1
			} else if i < 52 && j > 51 {
				cell.quadrant = 2
			} else if i > 52 && j < 51 {
				cell.quadrant = 3
			} else if i > 52 && j > 51 {
				cell.quadrant = 4
			}
			row = append(row, cell)
		}
		area = append(area, row)
	}
	return area
}
