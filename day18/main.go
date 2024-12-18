package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

const (
	charWall  = '#'
	charEmpty = '.'
)

type Direction byte

const (
	Up    Direction = '^'
	Down  Direction = 'v'
	Left  Direction = '<'
	Right Direction = '>'
)

type Point struct {
	x, y int
}

type State struct {
	pos   Point
	dir   Direction
	steps int
}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/18/input")

	fmt.Println("Part 1:")
	part1(data)

	fmt.Println("\nPart 2:")
	part2(data)
}

func part1(data []byte) {
	grid, walls := parse(data)

	// kilobyte of walls
	wallCount := 1024

	for i := 0; i < wallCount; i++ {
		wall := walls[i]
		grid[wall[1]][wall[0]] = charWall
	}

	steps := findPath(grid)
	if steps == -1 {
		fmt.Println("No path found!")
	} else {
		fmt.Printf("Path found in %d steps\n", steps)
	}
}

func part2(data []byte) {
	walls := parseWalls(data)
	width := 71
	height := 71

	for i := 0; i < len(walls); i++ {
		// fresh grid for each attempt
		grid := makeGrid(width, height)

		for j := 0; j <= i; j++ {
			wall := walls[j]
			grid[wall[1]][wall[0]] = charWall
		}
		if !hasPath(grid) {
			wall := walls[i]
			fmt.Printf("%d,%d\n", wall[0], wall[1])
			break
		}
	}
}

func findPath(grid [][]byte) int {
	queue := make([]State, 0, 1000)
	queue = append(queue, State{
		pos:   Point{0, 0},
		dir:   Right,
		steps: 0,
	})

	visited := make(map[Point]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.pos] {
			continue
		}
		visited[current.pos] = true

		if current.pos.x == 70 && current.pos.y == 70 {
			return current.steps
		}

		for _, dir := range []Direction{Up, Right, Down, Left} {
			nextPos := getNextPosition(current.pos, dir)
			if isValidMove(grid, nextPos) {
				queue = append(queue, State{
					pos:   nextPos,
					dir:   dir,
					steps: current.steps + 1,
				})
			}
		}

		sort.Slice(queue, func(i, j int) bool {
			return queue[i].steps < queue[j].steps
		})
	}

	return -1
}

func hasPath(grid [][]byte) bool {
	queue := make([]State, 0, 1000)
	queue = append(queue, State{
		pos:   Point{0, 0},
		dir:   Right,
		steps: 0,
	})

	visited := make(map[Point]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.pos] {
			continue
		}
		visited[current.pos] = true

		if current.pos.x == 70 && current.pos.y == 70 {
			return true
		}

		for _, dir := range []Direction{Up, Right, Down, Left} {
			nextPos := getNextPosition(current.pos, dir)
			if isValidMove(grid, nextPos) {
				queue = append(queue, State{
					pos:   nextPos,
					dir:   dir,
					steps: current.steps + 1,
				})
			}
		}

		sort.Slice(queue, func(i, j int) bool {
			return queue[i].steps < queue[j].steps
		})
	}

	return false
}

func getNextPosition(pos Point, dir Direction) Point {
	switch dir {
	case Up:
		return Point{pos.x, pos.y - 1}
	case Down:
		return Point{pos.x, pos.y + 1}
	case Left:
		return Point{pos.x - 1, pos.y}
	case Right:
		return Point{pos.x + 1, pos.y}
	}
	return pos
}

func isValidMove(grid [][]byte, p Point) bool {
	return p.x >= 0 && p.y >= 0 &&
		p.y < len(grid) && p.x < len(grid[0]) &&
		grid[p.y][p.x] == charEmpty
}

func makeGrid(width, height int) [][]byte {
	grid := make([][]byte, height)
	for j := 0; j < height; j++ {
		grid[j] = make([]byte, width)
		for k := 0; k < width; k++ {
			grid[j][k] = charEmpty
		}
	}
	return grid
}

func parse(data []byte) ([][]byte, [][2]int) {
	walls := parseWalls(data)
	grid := makeGrid(71, 71)
	return grid, walls
}

func parseWalls(data []byte) [][2]int {
	rows := bytes.Split(data, []byte("\n"))
	walls := make([][2]int, 0, len(rows))

	for i := range rows {
		parts := bytes.Split(rows[i], []byte(","))
		if len(parts) != 2 {
			continue
		}

		x, _ := strconv.Atoi(string(parts[0]))
		y, _ := strconv.Atoi(string(parts[1]))
		wall := [2]int{x, y}
		walls = append(walls, wall)
	}

	return walls
}
