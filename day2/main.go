package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	body := input.GetInput("https://adventofcode.com/2024/day/2/input")
	fmt.Println(handleInput(body))
}

func handleInput(body []byte) int {
	count := 0
	reports := bytes.Split(body, []byte("\n"))
	for i := 0; i < len(reports); i++ {
		levels := strings.Fields(string(reports[i]))
		if len(levels) == 0 {
			continue
		}

		broke := false
		strike := 0

		lastLevel, Err := strconv.Atoi(levels[0])
		if Err != nil {
			log.Fatal(Err)
		}

		var dir int
		firstLevel, err1 := strconv.Atoi(levels[0])
		finalLevel, err2 := strconv.Atoi(levels[len(levels)-1])

		secondLastLevel, err3 := strconv.Atoi(levels[len(levels)-2])
		secondLevel, err5 := strconv.Atoi(levels[1])

		if err1 != nil || err2 != nil || err3 != nil || err5 != nil {
			log.Fatal(err1)
		}
		if firstLevel > finalLevel && firstLevel > secondLastLevel {
			dir = -1
		} else if firstLevel < finalLevel && firstLevel < secondLastLevel {
			dir = 1
		}

		if dir == 0 {
			if firstLevel > secondLevel {
				dir = -1
			} else if firstLevel < secondLevel {
				dir = 1
			} else if dir == 0 {
				if firstLevel > finalLevel || firstLevel > secondLastLevel {
					dir = -1
				} else if firstLevel < finalLevel || firstLevel < secondLastLevel {
					dir = 1
				}
			}
		}

		for j := 1; j < len(levels); j++ {
			level, convErr := strconv.Atoi(levels[j])
			if convErr != nil {
				log.Fatal(convErr)
			}
			if (dir == -1 && level >= lastLevel-3 && level < lastLevel) || (dir == 1 && level <= lastLevel+3 && level > lastLevel) {
				lastLevel = level
				continue
			} else {
				strike++
				if strike < 2 {
					continue
				} else {
					broke = true
					break
				}
			}
		}

		if !broke {
			count++
		}
	}
	return count
}
