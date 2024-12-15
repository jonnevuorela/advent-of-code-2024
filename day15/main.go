package main

import (
	"aoc/input"
	"bytes"
	"fmt"
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

type Warehouse struct {
	grid   [][]byte
	player Point
}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/15/input")
	area, actions := parseInputSections(data)

	warehouse := parseWarehouseGrid(area)

	for _, action := range actions {
		warehouse.move(Direction(action))
	}

	result := warehouse.calculateGPS()
	fmt.Printf("Sum of GPS coordinates: %d\n", result)
}

func (w *Warehouse) move(dir Direction) {
	newPos := w.getNextPosition(dir)
	if !w.isInBounds(newPos) || w.grid[newPos.y][newPos.x] == '#' {
		return
	}

	if w.grid[newPos.y][newPos.x] == 'O' {
		boxNextPos := Point{
			x: newPos.x,
			y: newPos.y,
		}

		switch dir {
		case Up:
			boxNextPos.y--
		case Down:
			boxNextPos.y++
		case Left:
			boxNextPos.x--
		case Right:
			boxNextPos.x++
		}

		if !w.isInBounds(boxNextPos) || w.grid[boxNextPos.y][boxNextPos.x] != '.' {
			return
		}

		w.grid[boxNextPos.y][boxNextPos.x] = 'O'
		w.grid[newPos.y][newPos.x] = '.'
	}

	w.grid[w.player.y][w.player.x] = '.'
	w.player = newPos
	w.grid[w.player.y][w.player.x] = '@'
}

func (w *Warehouse) getNextPosition(dir Direction) Point {
	pos := w.player
	switch dir {
	case Up:
		pos.y--
	case Down:
		pos.y++
	case Left:
		pos.x--
	case Right:
		pos.x++
	}
	return pos
}

func (w *Warehouse) isInBounds(p Point) bool {
	return p.x >= 0 && p.y >= 0 && p.y < len(w.grid) && p.x < len(w.grid[0])
}

func (w *Warehouse) calculateGPS() int {
	sum := 0
	for y := range w.grid {
		for x := range w.grid[y] {
			if w.grid[y][x] == 'O' {
				sum += 100*y + x
			}
		}
	}
	return sum
}

func parseWarehouseGrid(data []byte) *Warehouse {
	if len(data) == 0 {
		return nil
	}

	rows := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	if len(rows) == 0 {
		return nil
	}

	w := &Warehouse{
		grid: make([][]byte, len(rows)),
	}

	for y, row := range rows {
		w.grid[y] = make([]byte, len(row))
		copy(w.grid[y], row)
		if x := bytes.IndexByte(row, '@'); x >= 0 {
			w.player = Point{x: x, y: y}
		}
	}
	return w
}

func parseInputSections(data []byte) ([]byte, []byte) {
	parts := bytes.Split(bytes.TrimSpace(data), []byte("\n\n"))

	area := parts[0]
	actions := bytes.Join(bytes.Fields(parts[1]), []byte{})

	return area, actions
}
