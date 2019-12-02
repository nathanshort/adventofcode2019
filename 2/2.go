package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func runProgram(input []int, noun int, verb int) []int {

	program := make([]int, len(input))
	copy(program, input)
	program[1] = noun
	program[2] = verb
	pc := 0

Forever:
	for {
		switch program[pc] {
		case 1:
			program[program[pc+3]] = program[program[pc+1]] + program[program[pc+2]]
		case 2:
			program[program[pc+3]] = program[program[pc+1]] * program[program[pc+2]]
		case 99:
			break Forever
		default:
			log.Fatalf("unknown opcode pc(%d) program(%v)", pc, program)
		}
		pc += 4
	}

	return program
}

func part1(program []int) {

	result := runProgram(program, 12, 2)
	fmt.Printf("part 1: %d\n", result[0])

}

func part2(program []int) {

	target := 19690720
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			result := runProgram(program, noun, verb)
			if result[0] == target {
				fmt.Printf("part 2: %d\n", 100*noun+verb)
				return
			}
		}

	}
	log.Fatal("should not get here")
}

func main() {

	input := "1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,9,1,19,1,19,5,23,1,23,6,27,2,9,27,31,1,5,31,35,1,35,10,39,1,39,10,43,2,43,9,47,1,6,47,51,2,51,6,55,1,5,55,59,2,59,10,63,1,9,63,67,1,9,67,71,2,71,6,75,1,5,75,79,1,5,79,83,1,9,83,87,2,87,10,91,2,10,91,95,1,95,9,99,2,99,9,103,2,10,103,107,2,9,107,111,1,111,5,115,1,115,2,119,1,119,6,0,99,2,0,14,0"
	var program []int
	for _, value := range strings.Split(input, ",") {
		asInt, _ := strconv.Atoi(value)
		program = append(program, asInt)
	}

	part1(program)
	part2(program)

}
