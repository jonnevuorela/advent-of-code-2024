package simulation

import (
	"aoc/types"
	"aoc/ui"
	"fmt"
	"time"
)

type UpdateMessage struct {
	Message1 string
	Message2 string
}
type Simulation struct {
	state     *types.State
	updateCh  chan [][]types.Tile
	messageCh chan UpdateMessage
	stopCh    chan struct{}
}

/**
* @brief Run simulation. Start as go routine
*
* @param
 */
func (s *Simulation) Run() {
	ticker := time.Tick(time.Second / 100)
	running := true

	for running {
		select {
		case <-ticker:
			s.state.Buffer.Reset()
			_, rows := ui.GetTerminalSize()
			s.state.ScaleFactor = float64(rows) * 0.9 / float64(len(s.state.Data))
			ui.UpdateDisplay(s.state)

		case newData := <-s.updateCh:
			s.state.Data = newData

		case newMessages := <-s.messageCh:
			s.state.Message1 = newMessages.Message1
			s.state.Message2 = newMessages.Message2

		case <-s.stopCh:
			running = false
		}
	}
}

/**
* @brief Setup a new simulation with starting data, and establish needed channels for updating simulation.
*
* @param [][]types.Tile
*
* @return
 */
func NewSimulation(initialData [][]types.Tile) *Simulation {
	//clear terminal
	fmt.Print("\033[2J")
	return &Simulation{
		state:     NewState(initialData),
		updateCh:  make(chan [][]types.Tile, 1),
		messageCh: make(chan UpdateMessage, 1),
		stopCh:    make(chan struct{}),
	}
}

/**
* @brief Send new data to simulation trought established channel.
*
* @param
 */
func (s *Simulation) UpdateData(newData [][]types.Tile) {
	select {
	case s.updateCh <- newData:
	default:
		// Channel is full, skip this update
	}
}

/**
* @brief Send new messages to simulation through established channel.
*
* @param msg1 string First message
* @param msg2 string Second message
 */
func (s *Simulation) UpdateMessages(msg1, msg2 string) {
	select {
	case s.messageCh <- UpdateMessage{Message1: msg1, Message2: msg2}:
	default:
		// Channel is full, skip this update
	}
}

/**
* @brief stops the simulation.
*
* @param
 */
func (s *Simulation) Stop() {
	close(s.stopCh)
}

/**
* @brief initialize new state with optional messages
*
* @param data [][]types.Tile
* @param messages ...string - optional messages (up to 2)
*
* @return *types.State
 */
func NewState(data [][]types.Tile, messages ...string) *types.State {
	lastPrinted := make([][]byte, len(data))
	for i := range lastPrinted {
		lastPrinted[i] = make([]byte, len(data[i]))
	}

	// Default messages
	msg1 := ""
	msg2 := ""

	// Override with provided messages if any
	if len(messages) > 0 {
		msg1 = messages[0]
		if len(messages) > 1 {
			msg2 = messages[1]
		}
	}

	return &types.State{
		Data:        data,
		LastPrinted: lastPrinted,
		Message1:    msg1,
		Message2:    msg2,
	}
}
