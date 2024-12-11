package ui

import (
	"aoc/types"
	"bytes"
	"fmt"
	"sync"
)

type HighlightType int

const (
	OldHighlight HighlightType = iota
	NewHighlight
)

type HighlightInfo struct {
	isHighlighted bool
	isNew         bool
}

var (
	highlightedPositions = make(map[[2]int]HighlightInfo)
	highlightMutex       sync.RWMutex
)

func UpdateDisplay(state *types.State) {
	updateDataBuffer(&state.Buffer, state.Data, state.LastPrinted, state.ScaleFactor)
	drawBorder(&state.Buffer, len(state.Data), len(state.Data[0])*2, state.ScaleFactor)
	writeMessages(&state.Buffer, state.Message1, state.Message2)
	renderBuffer(state.Buffer.String())
}

func writeMessages(buffer *bytes.Buffer, message1, message2 string) {
	_, rows := GetTerminalSize()
	buffer.WriteString(fmt.Sprintf("\033[%d;1H", rows-2))
	buffer.WriteString(message1)
	buffer.WriteString(fmt.Sprintf("\033[%d;1H", rows-1))
	buffer.WriteString(message2)
}

func renderBuffer(content string) {
	fmt.Print("\033[H")
	fmt.Print(content)
}

func SetHighlighted(pos [2]int) {
	highlightMutex.Lock()
	highlightedPositions[pos] = HighlightInfo{isHighlighted: true, isNew: true}
	highlightMutex.Unlock()
}

func ClearHighlights() {
	highlightMutex.Lock()
	highlightedPositions = make(map[[2]int]HighlightInfo)
	highlightMutex.Unlock()
}

func RemoveHighlight(pos [2]int) {
	highlightMutex.Lock()
	delete(highlightedPositions, pos)
	highlightMutex.Unlock()
}
func ConvertNewToOldHighlights() {
	highlightMutex.Lock()
	for pos, info := range highlightedPositions {
		if info.isNew {
			highlightedPositions[pos] = HighlightInfo{isHighlighted: true, isNew: false}
		}
	}
	highlightMutex.Unlock()
}

func GetHighlightedPositions() map[[2]int]HighlightInfo {
	highlightMutex.RLock()
	defer highlightMutex.RUnlock()

	// Return a copy to be safe
	result := make(map[[2]int]HighlightInfo)
	for k, v := range highlightedPositions {
		result[k] = v
	}
	return result
}

func updateDataBuffer(buffer *bytes.Buffer, cell [][]types.Tile, lastPrinted [][]byte, scaleFactor float64) {
	var dataBuffer bytes.Buffer
	dataColor := GetColor("red")
	oldHighlightColor := GetColor("bgYellow")
	newHighlightColor := GetColor("bgGreen")
	reset := "\033[0m"

	highlightMutex.RLock()
	defer highlightMutex.RUnlock()

	for i := range cell {
		dataBuffer.WriteString(fmt.Sprintf("\033[%d;%dH", int(float64(i+2)*scaleFactor), 2))
		for j := range cell[i] {
			if highlight, exists := highlightedPositions[cell[i][j].Location]; exists {
				if highlight.isNew {
					dataBuffer.WriteString(dataColor + newHighlightColor)
				} else {
					dataBuffer.WriteString(dataColor + oldHighlightColor)
				}
			} else {
				dataBuffer.WriteString(dataColor)
			}

			dataBuffer.WriteString(fmt.Sprintf("%d ", cell[i][j].Value))
			dataBuffer.WriteString(reset)
			lastPrinted[i][j] = byte(cell[i][j].Value)
		}
	}
	buffer.Write(dataBuffer.Bytes())
}

func writeRow(buffer *bytes.Buffer, cellRow []types.Tile, lastRow []byte, rowIndex int, scaleFactor float64, color, reset string) {
	buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s", int(float64(rowIndex+2)*scaleFactor), 2, color))
	for j, tile := range cellRow {
		buffer.WriteString(fmt.Sprintf("%d ", tile.Value))
		lastRow[j] = byte(tile.Value)
	}
	buffer.WriteString(reset)
}
