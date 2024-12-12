package main

import (
	"aoc/input"
	"bytes"
	"fmt"
)

type region struct {
	plantType byte
	area      int
	per       int
}
type plant struct {
	plant    byte
	location [2]int
	region   region
}

var directions = [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/12/input")
	plants, grid := mapPlants(data)
	if len(grid) == 0 || len(grid[0]) == 0 {
		fmt.Println("Error: Empty grid")
		return
	}
	regions := mapRegions(plants, grid)

	totalPrice := 0
	for _, r := range regions {
		price := r.area * r.per
		totalPrice += price
	}

	fmt.Println("Total price:", totalPrice)
}

func mapRegions(plants []plant, grid [][]byte) []region {
	regions := []region{}
	visited := make([][]bool, len(grid))
	for i := range visited {
		visited[i] = make([]bool, len(grid[0]))
	}

	for i := range plants {
		if plants[i].region != (region{}) {
			continue
		}
		area, perimeter := mapSingleRegion(plants[i], grid, visited)
		newRegion := region{
			plantType: plants[i].plant,
			area:      area,
			per:       perimeter,
		}
		regions = append(regions, newRegion)
		plants[i].region = newRegion
	}
	return regions
}

func mapSingleRegion(start plant, grid [][]byte, visited [][]bool) (int, int) {
	currentRow := start.location[1]
	currentCol := start.location[0]

	if visited[currentRow][currentCol] {
		return 0, 0
	}

	visited[currentRow][currentCol] = true
	area := 1
	perimeter := 0

	// check all sides of current cell
	for _, dir := range directions {
		newRow := currentRow + dir[0]
		newCol := currentCol + dir[1]

		if newRow < 0 || newRow >= len(grid) ||
			newCol < 0 || newCol >= len(grid[0]) ||
			grid[newRow][newCol] != start.plant {
			perimeter++
			continue
		}

		if !visited[newRow][newCol] {
			nextPlant := plant{
				plant:    start.plant,
				location: [2]int{newCol, newRow},
				region:   start.region,
			}
			a, p := mapSingleRegion(nextPlant, grid, visited)
			area += a
			perimeter += p
		}
	}
	return area, perimeter
}

func mapPlants(data []byte) ([]plant, [][]byte) {
	plants := []plant{}
	rows := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	if len(rows) == 0 {
		return nil, nil
	}

	grid := make([][]byte, len(rows))

	for i := range rows {
		grid[i] = make([]byte, len(rows[i]))
		for j := range rows[i] {
			grid[i][j] = rows[i][j]
			np := plant{
				plant:    rows[i][j],
				location: [2]int{j, i},
			}
			plants = append(plants, np)
		}
	}

	return plants, grid
}
