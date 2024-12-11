package main

import (
	"aoc/input"
	"aoc/simulation"
	"aoc/types"
	"aoc/ui"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Cell struct {
	types.Tile
	isTrailhead bool
}

type TrailStatus int

type trail struct {
	trailhead Cell
	trail     []Cell
	status    TrailStatus
}

const (
	Active TrailStatus = iota
	DNF
	Completed
)

var directions = [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/10/input")
	md := mapInput(data)

	cells := make([][]Cell, len(md))
	for i := range md {
		cells[i] = make([]Cell, len(md[i]))
		for j := range md[i] {
			cells[i][j] = Cell{
				Tile:        md[i][j],
				isTrailhead: false,
			}
		}
	}

	sim := simulation.NewSimulation(md)
	go sim.Run()

	var activeTrails []trail
	var completedTrails []trail
	reachableNines := make(map[[2]int]map[[2]int]bool)

	go func() {
		displayTicker := time.NewTicker(100 * time.Millisecond)
		progressTicker := time.NewTicker(1 * time.Second)
		defer displayTicker.Stop()
		defer progressTicker.Stop()

		checkTrailheads(cells)
		activeTrails = formTrails(cells)
		for _, t := range activeTrails {
			reachableNines[t.trailhead.Location] = make(map[[2]int]bool)
		}

		for {
			select {
			case <-displayTicker.C:
				sim.UpdateData(updateData(cells))
				totalScore := 0
				for _, nines := range reachableNines {
					totalScore += len(nines)
				}

				activeCount := 0
				for _, t := range activeTrails {
					if t.status == Active {
						activeCount++
					}
				}

				sim.UpdateMessages(
					fmt.Sprintf("Active: %d, Completed: %d, DNF: %d",
						activeCount, len(completedTrails), countDNFTrails(activeTrails)),
					fmt.Sprintf("Score: %d", totalScore),
				)

			case <-progressTicker.C:
				ui.ConvertNewToOldHighlights()
				if newCompleted := progress(activeTrails, cells, reachableNines); len(newCompleted) > 0 {
					completedTrails = append(completedTrails, newCompleted...)
				}
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	sim.Stop()
}

func checkTrailheads(data [][]Cell) {
	for i := range data {
		for j := range data[i] {
			if data[i][j].Value == 0 {
				data[i][j].isTrailhead = true
				ui.SetHighlighted(data[i][j].Location)
			}
		}
	}
}
func newTrail(th Cell) trail {
	return trail{
		trailhead: th,
		trail:     []Cell{th},
	}
}

func progress(trails []trail, data [][]Cell, reachableNines map[[2]int]map[[2]int]bool) []trail {
	var completedTrails []trail
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for i := range trails {
		if trails[i].status != Active {
			continue
		}

		currentCell := trails[i].trail[len(trails[i].trail)-1]
		if currentCell.Value == 9 {
			trails[i].status = Completed
			completedTrails = append(completedTrails, trails[i])
			reachableNines[trails[i].trailhead.Location][currentCell.Location] = true
			continue
		}

		found := false
		row, col := currentCell.Location[1]-1, currentCell.Location[0]-1

		for _, dir := range directions {
			if found {
				break
			}

			newRow, newCol := row+dir[0], col+dir[1]
			if newRow < 0 || newRow >= len(data) || newCol < 0 || newCol >= len(data[0]) {
				continue
			}

			nextCell := data[newRow][newCol]
			if nextCell.Value != currentCell.Value+1 {
				continue
			}

			visited := false
			for _, cell := range trails[i].trail {
				if cell.Location == nextCell.Location {
					visited = true
					break
				}
			}

			if !visited {
				trails[i].trail = append(trails[i].trail, nextCell)
				ui.SetHighlighted(nextCell.Location)
				found = true
			}
		}

		if !found {
			trails[i].status = DNF
			for _, cell := range trails[i].trail {
				ui.RemoveHighlight(cell.Location)
			}
		}
	}

	return completedTrails
}
func formTrails(data [][]Cell) []trail {
	trails := []trail{}
	for i := range data {
		for j := range data[i] {
			if !data[i][j].isTrailhead {
				continue
			}

			for _, dir := range directions {
				newRow, newCol := i+dir[0], j+dir[1]
				if newRow < 0 || newRow >= len(data) || newCol < 0 || newCol >= len(data[0]) {
					continue
				}

				if isValidNextCell(data[i][j], data[newRow][newCol]) {
					ui.SetHighlighted(data[newRow][newCol].Location)
					t := newTrail(data[i][j])
					t.trail = append(t.trail, data[newRow][newCol])
					trails = append(trails, t)
				}
			}
		}
	}
	return trails
}
func countDNFTrails(trails []trail) int {
	count := 0
	for _, t := range trails {
		if t.status == DNF {
			count++
		}
	}
	return count
}

func isValidNextCell(currentCell, nextCell Cell) bool {
	return nextCell.Value == currentCell.Value+1
}

func updateData(data [][]Cell) [][]types.Tile {
	newData := make([][]types.Tile, len(data))
	for i := range data {
		newData[i] = make([]types.Tile, len(data[i]))
		for j := range data[i] {
			newData[i][j] = data[i][j].Tile
		}
	}
	return newData
}

func mapInput(data []byte) [][]types.Tile {
	var tiles [][]types.Tile
	rows := bytes.Split(data, []byte("\n"))
	for i, row := range rows {
		if len(row) == 0 {
			continue
		}
		var tileRow []types.Tile
		for j := range row {
			intVal, err := strconv.Atoi(string(row[j]))
			if err != nil {
				fmt.Print(err)
				continue
			}
			t := types.Tile{
				Value:    intVal,
				Location: [2]int{j + 1, i + 1},
			}
			tileRow = append(tileRow, t)
		}
		tiles = append(tiles, tileRow)
	}
	return tiles
}
