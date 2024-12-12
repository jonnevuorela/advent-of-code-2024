package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type stone struct {
	num []byte
}

func main() {
	startTime := time.Now()

	data := input.GetInput("https://adventofcode.com/2024/day/11/input")
	nums := parseDataInt(data)
	totalIterations := 75

	for i := 0; i < totalIterations; i++ {
		percentage := (float64(i) / float64(totalIterations)) * 100
		fmt.Printf("\rProgress: %.1f%%", percentage)
		newNums := blinkInt(nums)
		nums = newNums
		newNums = nil
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\rProgress: 100.0%%\n")
	fmt.Printf("Result: %d\n", len(nums))
	fmt.Printf("Execution time: %v\n", elapsed)
}

func blink(stones []stone) []stone {
	newStones := make([]stone, 0, len(stones)*2)

	for i := range stones {
		stones[i].num = bytes.TrimSpace(stones[i].num)

		if bytes.Equal(stones[i].num, []byte("0")) {
			stones[i].num = []byte("1")
			newStones = append(newStones, stones[i])
			continue
		}

		if len(stones[i].num)%2 == 0 {
			mid := len(stones[i].num) / 2
			firstHalf := stones[i].num[:mid]
			secondHalf := stones[i].num[mid:]

			for len(firstHalf) > 0 {
				firstHalf = bytes.TrimLeft(firstHalf, "0")
				if len(firstHalf) > 0 && firstHalf[0] == '0' {
					continue
				}
				break
			}
			if len(firstHalf) == 0 {
				firstHalf = []byte("0")
			}

			for len(secondHalf) > 0 {
				secondHalf = bytes.TrimLeft(secondHalf, "0")
				if len(secondHalf) > 0 && secondHalf[0] == '0' {
					continue
				}
				break
			}
			if len(secondHalf) == 0 {
				secondHalf = []byte("0")
			}

			stones[i].num = firstHalf
			newStones = append(newStones, stones[i])
			newStones = append(newStones, stone{num: secondHalf})
			continue
		}

		val, err := strconv.ParseInt(string(stones[i].num), 10, 64)
		if err != nil {
			fmt.Print(err)
			continue
		}
		stones[i].num = []byte(strconv.FormatInt(val*2024, 10))
		newStones = append(newStones, stones[i])
	}

	return newStones
}

func parseData(data []byte) []stone {
	elems := bytes.Split(data, []byte(" "))
	stones := []stone{}
	for i := range elems {
		stone := stone{
			num: elems[i],
		}
		stones = append(stones, stone)
	}

	return stones
}

func blinkInt(nums []int64) []int64 {
	newNums := make([]int64, 0, len(nums)*2)

	for _, num := range nums {
		if num == 0 {
			newNums = append(newNums, 1)
			continue
		}

		digits := 1
		for n := num; n >= 10; n /= 10 {
			digits++
		}

		if digits%2 == 0 {
			divisor := int64(1)
			for i := 0; i < digits/2; i++ {
				divisor *= 10
			}
			firstHalf := num / divisor
			secondHalf := num % divisor

			if firstHalf == 0 {
				firstHalf = 0
			}
			if secondHalf == 0 {
				secondHalf = 0
			}

			newNums = append(newNums, firstHalf)
			newNums = append(newNums, secondHalf)
		} else {
			newNums = append(newNums, num*2024)
		}
	}

	return newNums
}

func parseDataInt(data []byte) []int64 {
	elems := bytes.Split(data, []byte(" "))
	nums := make([]int64, 0, len(elems))

	for _, elem := range elems {
		elem = bytes.TrimSpace(elem)
		if len(elem) == 0 {
			continue
		}

		num, err := strconv.ParseInt(string(elem), 10, 64)
		if err != nil {
			fmt.Printf("Error parsing number: %v\n", err)
			continue
		}
		nums = append(nums, num)
	}

	return nums
}
