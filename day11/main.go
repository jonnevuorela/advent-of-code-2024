package main

import (
	"aoc/input"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync"
	"time"
)

type stone struct {
	num []byte
}

var blinkIntCache sync.Map

func main() {
	go func() {
		fmt.Println("Starting pprof server on :6060")
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	startTime := time.Now()

	data := input.GetInput("https://adventofcode.com/2024/day/11/input")
	nums := parseDataInt(data)
	totalIterations := 75

	for i := 0; i < totalIterations; i++ {
		if i < 45 {
			if i%5 == 0 {
				cleanCache()
			}
		} else if i >= 45 {
			cleanCache()
		}

		percentage := (float64(i) / float64(totalIterations)) * 100
		fmt.Printf("\rProgress: %.1f%%", percentage)
		nums = blinkInt(nums)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\rProgress: 100.0%%\n")
	fmt.Printf("Result: %d\n", len(nums))
	fmt.Printf("Execution time: %v\n", elapsed)

	var count int
	blinkIntCache.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	fmt.Printf("Cache size: %d\n", count)
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func createCacheKey(nums []int64) string {
	h := fnv.New64()
	for _, num := range nums {
		binary.Write(h, binary.LittleEndian, num)
	}
	return string(h.Sum(nil))
}

func cleanCache() {
	blinkIntCache = sync.Map{}
}

func countDigits(n int64) int {
	if n == 0 {
		return 1
	}
	return int(math.Log10(float64(n))) + 1
}

func blinkInt(nums []int64) []int64 {
	chunkSize := 100
	numChunks := (len(nums) + chunkSize - 1) / chunkSize

	// allocate exact size of needed
	exactCap := 0
	for _, num := range nums {
		if num == 0 || countDigits(num)%2 == 0 {
			exactCap += 2
		} else {
			exactCap += 1
		}
	}
	newNums := make([]int64, 0, exactCap)

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(nums) {
			end = len(nums)
		}
		chunk := nums[start:end]

		chunkKey := createCacheKey(chunk)
		if cachedChunk, exists := blinkIntCache.Load(chunkKey); exists {
			newNums = append(newNums, cachedChunk.([]int64)...)
			continue
		}

		chunkResult := processChunk(chunk)
		blinkIntCache.Store(chunkKey, chunkResult)
		newNums = append(newNums, chunkResult...)
	}

	return newNums
}
func willOverflow(num int64) bool {
	maxInt64 := int64(9223372036854775807)
	if num > maxInt64/2024 {
		return true
	}
	return false
}
func processChunk(chunk []int64) []int64 {
	result := make([]int64, 0, len(chunk)*2)
	for _, num := range chunk {
		if countDigits(num)%2 != 0 {
			if willOverflow(num) {
				panic("Number would overflow int64")
			}
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

			result = append(result, firstHalf)
			result = append(result, secondHalf)
		} else {
			result = append(result, num*2024)
		}
	}

	return result
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

// part 1
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
