package main

import (
	"aoc/input"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type op int

const (
	adv op = iota
	bxl
	bst
	jnz
	bxc
	out
	bdv
	cdv
)

type register struct {
	name  byte
	value int
}

type instruction struct {
	opcode  op
	operand int
}

func main() {
	data := input.GetInput("https://adventofcode.com/2024/day/17/input")
	regs, instructions := parse(data)
	output := []string{}
	instructionPointer := 0

	for instructionPointer < len(instructions)*2 {
		inst := instructions[instructionPointer/2]
		out := execute(inst, &regs, &instructionPointer)
		if out != -1 {
			output = append(output, fmt.Sprintf("%d", out))
		}
	}

	fmt.Println(strings.Join(output, ","))
}

func getComboValue(operand int, regs *[]register) int {
	if operand <= 3 {
		return operand
	}

	if operand == 4 {
		return (*regs)[0].value // A
	}
	if operand == 5 {
		return (*regs)[1].value // B
	}
	if operand == 6 {
		return (*regs)[2].value // C
	}
	return 0 // operand 7 is reserved
}

func execute(inst instruction, regs *[]register, instructionPointer *int) int {
	if inst.opcode < 0 || inst.opcode > 7 {
		fmt.Println("no such operation")
		return -1
	}

	var regA, regB, regC *register
	for i := range *regs {
		switch (*regs)[i].name {
		case 'A':
			regA = &(*regs)[i]
		case 'B':
			regB = &(*regs)[i]
		case 'C':
			regC = &(*regs)[i]
		}
	}

	switch inst.opcode {
	case adv:
		regA.value = regA.value / (1 << getComboValue(inst.operand, regs))
		*instructionPointer += 2
	case bxl:
		regB.value ^= inst.operand
		*instructionPointer += 2
	case bst:
		regB.value = getComboValue(inst.operand, regs) % 8
		*instructionPointer += 2
	case jnz:
		if regA.value != 0 {
			*instructionPointer = inst.operand * 2
		} else {
			*instructionPointer += 2
		}
	case bxc:
		regB.value = regB.value ^ regC.value
		*instructionPointer += 2
	case out:
		out := getComboValue(inst.operand, regs) % 8
		*instructionPointer += 2
		return out
	case bdv:
		regB.value = regA.value / (1 << getComboValue(inst.operand, regs))
		*instructionPointer += 2
	case cdv:
		regC.value = regA.value / (1 << getComboValue(inst.operand, regs))
		*instructionPointer += 2
	}
	return -1
}

func parse(data []byte) ([]register, []instruction) {
	parts := bytes.SplitN(bytes.TrimSpace(data), []byte("\n\nProgram: "), 2)

	registers := make([]register, 0, 3)
	regLines := bytes.Split(parts[0], []byte("\n"))
	for _, line := range regLines {
		parts := bytes.SplitN(bytes.TrimPrefix(line, []byte("Register ")), []byte(": "), 2)
		val, _ := strconv.Atoi(string(parts[1]))
		registers = append(registers, register{
			name:  parts[0][0],
			value: val,
		})
	}

	program := strings.TrimSpace(string(parts[1]))
	numStrs := strings.Split(program, ",")
	instructions := make([]instruction, 0, len(numStrs)/2)

	for i := 0; i < len(numStrs); i += 2 {
		if i+1 >= len(numStrs) {
			break
		}
		opcode, _ := strconv.Atoi(strings.TrimSpace(numStrs[i]))
		operand, _ := strconv.Atoi(strings.TrimSpace(numStrs[i+1]))
		instructions = append(instructions, instruction{
			opcode:  op(opcode),
			operand: operand,
		})
	}

	return registers, instructions
}
