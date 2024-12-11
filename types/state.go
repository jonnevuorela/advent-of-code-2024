package types

import "bytes"

type Tile struct {
	Value    int
	Location [2]int
}

type State struct {
	Data        [][]Tile
	LastPrinted [][]byte
	Message1    string
	Message2    string
	Buffer      bytes.Buffer
	ScaleFactor float64
}
