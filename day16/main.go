package main

import (
	"aoc/input"
	"bytes"
	"fmt"
)

const (
	initialQueueCapacity = 1000
	initialMapCapacity   = 10000

	moveCost = 1
	turnCost = 1000

	charStart = 'S'
	charEnd   = 'E'
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

func (p Point) IsZero() bool {
	return p.x == 0 && p.y == 0
}

type ReindeerState struct {
	pos   Point
	dir   Direction
	score int
}

type VisitedKey struct {
	pos Point
	dir Direction
}

type Maze struct {
	grid     [][]byte
	reindeer Point
}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/16/input")
	maze := mapInput(data)
	if maze == nil {
		return
	}

	bestScore, remaining := maze.ExploreReindeerPaths()

	if bestScore == int(^uint(0)>>1) {
		fmt.Println("No reindeer reached the end!")
	} else {
		fmt.Printf("Best reindeer path score: %d\n", bestScore)
	}

	if remaining > 0 {
		fmt.Printf("Unfinished routes: %d\n", remaining)
	} else {
		fmt.Println("All routes completed!")
	}
}

// Dijkstra with priority que
func (m *Maze) ExploreReindeerPaths() (int, int) {
	// priority que
	queue := make([]ReindeerState, 0, initialQueueCapacity)
	queue = append(queue, ReindeerState{
		pos:   m.reindeer,
		dir:   Right,
		score: 0,
	})

	visited := make(map[VisitedKey]int, initialMapCapacity)
	bestScore := int(^uint(0) >> 1) // worst case scenario as default
	completedPaths := 0

	for len(queue) > 0 {
		minIdx := 0
		// compare score to best in que
		for i := 1; i < len(queue); i++ {
			if queue[i].score < queue[minIdx].score {
				minIdx = i // switch if bested
			}
		}

		// take out best from que
		current := queue[minIdx]
		queue[minIdx] = queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if m.grid[current.pos.y][current.pos.x] == charEnd {
			completedPaths++
			if current.score < bestScore {
				bestScore = current.score // best so result so far
			}
			continue
		}

		// check if we've been here before with better score
		key := VisitedKey{pos: current.pos, dir: current.dir}
		if prevScore, exists := visited[key]; exists && prevScore <= current.score {
			continue
		}
		visited[key] = current.score

		// try moving forward in same direction
		if nextPos := m.getNextPosition(current.pos, current.dir); m.isValidMove(nextPos) {
			queue = append(queue, ReindeerState{
				pos:   nextPos,
				dir:   current.dir,
				score: current.score + moveCost,
			})
		}

		for _, newDir := range []Direction{
			current.dir.RotateClockwise(),
			current.dir.RotateCounterClockwise(),
		} {
			if nextPos := m.getNextPosition(current.pos, newDir); m.isValidMove(nextPos) {
				queue = append(queue, ReindeerState{
					pos:   nextPos,
					dir:   newDir,
					score: current.score + turnCost + moveCost,
				})
			}
		}
	}

	return bestScore, len(queue)
}

func (d Direction) RotateClockwise() Direction {
	switch d {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	default:
		return d
	}
}

func (d Direction) RotateCounterClockwise() Direction {
	switch d {
	case Up:
		return Left
	case Left:
		return Down
	case Down:
		return Right
	case Right:
		return Up
	default:
		return d
	}
}

func (m *Maze) getNextPosition(pos Point, dir Direction) Point {
	newPos := pos
	switch dir {
	case Up:
		newPos.y--
	case Down:
		newPos.y++
	case Left:
		newPos.x--
	case Right:
		newPos.x++
	}
	return newPos
}

func (m *Maze) isValidMove(p Point) bool {
	return m.isInBounds(p) && (m.grid[p.y][p.x] == charEmpty || m.grid[p.y][p.x] == charEnd)
}

func (m *Maze) isInBounds(p Point) bool {
	return p.x >= 0 && p.y >= 0 && p.y < len(m.grid) && p.x < len(m.grid[0])
}

func mapInput(data []byte) *Maze {
	if len(data) == 0 {
		return nil
	}

	rows := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	if len(rows) == 0 {
		return nil
	}

	m := &Maze{
		grid: make([][]byte, len(rows)),
	}

	for y, row := range rows {
		m.grid[y] = make([]byte, len(row))
		copy(m.grid[y], row)
		if x := bytes.IndexByte(row, charStart); x >= 0 {
			m.reindeer = Point{x: x, y: y}
			m.grid[y][x] = charEmpty
		}
	}

	if m.reindeer.IsZero() {
		return nil
	}

	return m
}
