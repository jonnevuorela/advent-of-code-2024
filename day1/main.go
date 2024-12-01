package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	input := getInput()
	left, right := handleInput(input)
	left = sortList(left)
	right = sortList(right)
	answer := calculateDifference(left, right)
	fmt.Println("Difference is ", answer)
	answer2 := similarityScore(left, right)
	fmt.Println("Similarity Score is ", answer2)
}

func getEnvVar(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv(key)
}
func similarityScore(left []int, right []int) int {
	var value int
	for i := 0; i < len(left); i++ {
		multiplier := 0
		for in := 0; in < len(right); in++ {
			if left[i] == right[in] {
				multiplier++
			}
		}
		val := left[i] * multiplier
		value += val
	}
	return value
}
func calculateDifference(left []int, right []int) int {
	var value int
	for i := 0; i < len(left); i++ {
		diff := math.Abs(float64(left[i] - right[i]))
		value += int(diff)
	}
	return value
}
func sortList(list []int) []int {
	for i := 0; i < len(list)-1; i++ {
		for index := 0; index < len(list)-1-i; index++ {
			if list[index] > list[index+1] {
				list[index] ^= list[index+1]
				list[index+1] ^= list[index]
				list[index] ^= list[index+1]
			}
		}
	}
	return list
}
func handleInput(body []byte) ([]int, []int) {
	elements := strings.Fields(string(body))
	var right []int
	var left []int
	for i := 0; i < len(elements); i++ {
		value, convErr := strconv.Atoi(elements[i])
		if convErr != nil {
			log.Fatal(convErr)
		}
		if i%2 == 0 {

			left = append(left, value)
		} else {
			right = append(right, value)
		}
	}
	return left, right
}

func getInput() []byte {
	url := "https://adventofcode.com/2024/day/1/input"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:   "session",
		Value:  getEnvVar("SESSION"),
		Path:   "/",
		Domain: "adventofcode.com",
	}

	req.AddCookie(cookie)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body
}
