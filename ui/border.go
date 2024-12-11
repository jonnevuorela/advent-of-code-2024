package ui

import (
	"bytes"
	"fmt"
)

type borderChars struct {
	topLeft     string
	topRight    string
	bottomLeft  string
	bottomRight string
	horizontal  string
	vertical    string
}

var defaultBorder = borderChars{
	topLeft:     "╭",
	topRight:    "╮",
	bottomLeft:  "╰",
	bottomRight: "╯",
	horizontal:  "─",
	vertical:    "│",
}

// drawBorder draws a border around the content
func drawBorder(buffer *bytes.Buffer, rows, cols int, scaleFactor float64) {
	color := GetColor("brightBlue")
	reset := "\033[0m"

	drawTopBorder(buffer, cols, color, reset)
	drawSideBorders(buffer, rows, cols, scaleFactor, color, reset)
	drawBottomBorder(buffer, rows, cols, scaleFactor, color, reset)
}

func drawTopBorder(buffer *bytes.Buffer, cols int, color, reset string) {
	buffer.WriteString("\033[1;1H")
	buffer.WriteString(color + defaultBorder.topLeft)
	for i := 0; i < cols; i++ {
		buffer.WriteString(defaultBorder.horizontal)
	}
	buffer.WriteString(defaultBorder.topRight + reset)
}

func drawSideBorders(buffer *bytes.Buffer, rows, cols int, scaleFactor float64, color, reset string) {
	for i := 0; i < int(float64(rows)*scaleFactor); i++ {
		buffer.WriteString(fmt.Sprintf("\033[%d;1H%s%s%s", i+2, color, defaultBorder.vertical, reset))
		buffer.WriteString(fmt.Sprintf("\033[%d;%dH%s%s%s", i+2, cols+2, color, defaultBorder.vertical, reset))
	}
}

func drawBottomBorder(buffer *bytes.Buffer, rows, cols int, scaleFactor float64, color, reset string) {
	buffer.WriteString(fmt.Sprintf("\033[%d;1H%s%s", int(float64(rows)*scaleFactor)+2, color, defaultBorder.bottomLeft))
	for i := 0; i < cols; i++ {
		buffer.WriteString(defaultBorder.horizontal)
	}
	buffer.WriteString(defaultBorder.bottomRight + reset)
}
