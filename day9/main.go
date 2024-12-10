package main

import (
	"aoc/input"
	"fmt"
	"strconv"
	"unicode"
)

type block struct {
	b      []byte
	id     int
	isData bool
}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/9/input")
	pd := defragment(expandData(expandFreeSpace(parseData(data))))

	setId(pd)
	fmt.Print(checkChecksum(pd))
}

func setId(data []block) {
	idCounter := 0
	for i := range data {
		if string(data[i].b) != "." {
			data[i].id = idCounter
			idCounter++
		}
	}
}

func checkChecksum(data []block) int {
	var cs int
	for i := range data {
		if string(data[i].b) != "." {
			num1, err := strconv.Atoi(string(data[i].b))
			num2 := data[i].id
			if err != nil {
				fmt.Print(err)
			}
			cs += num1 * num2
		}
	}
	return cs
}

func defragment(data []block) []block {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i].isData {
			s := checkSpace(data)
			if s != -1 && s < i {
				data[i], data[s] = data[s], data[i]
			}
		}
	}
	return data
}

func checkSpace(data []block) int {
	for i := range data {
		if !data[i].isData {
			return i
		}
	}
	return -1
}

func parseData(data []byte) []block {
	var pd []block
	for i := range data {
		if unicode.IsDigit(rune(data[i])) {
			if (i+1)%2 == 0 {
				fs := block{
					b:      []byte{data[i]},
					isData: false,
				}
				pd = append(pd, fs)
			} else {
				db := block{
					b:      []byte{data[i]},
					id:     countDataBlocks(pd),
					isData: true,
				}
				pd = append(pd, db)
			}
		}
	}
	return pd
}

func expandData(data []block) []block {
	var expd []block
	for i := range data {
		if data[i].isData {
			if data[i].b[0] != '\n' {
				l, err := strconv.Atoi(string(data[i].b))
				if err != nil {
					fmt.Print(err)
				}
				if l != 0 {
					data[i].b = []byte(strconv.Itoa(data[i].id))
					for j := 0; j < l; j++ {
						expd = append(expd,
							block{
								b:      []byte(strconv.Itoa(data[i].id)),
								isData: true,
							})
					}
				}
			}
		} else {
			expd = append(expd, data[i])
		}
	}
	return expd
}

func expandFreeSpace(data []block) []block {
	var expd []block
	for i := range data {
		if !data[i].isData {
			if data[i].b[0] != '\n' {
				l, err := strconv.Atoi(string(data[i].b))
				if err != nil {
					fmt.Print(err)
				}
				if l != 0 {
					data[i].b = []byte(".")
					for j := 0; j < l; j++ {
						expd = append(expd,
							block{
								b:      []byte("."),
								isData: false,
							})
					}
				}
			}
		} else {
			expd = append(expd, data[i])
		}
	}
	return expd
}

func countDataBlocks(data []block) int {
	num := 0
	for i := range data {
		if data[i].isData {
			num++
		}
	}
	return num
}
