package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func runProgram(input []int, inputValue int) {

	program := make([]int, len(input))
	copy(program, input)
	pc := 0

	for {
		opcode := program[pc]%10 + program[pc]/10%10*10
		params := make([]int, 2)

		switch opcode {
		case 1, 2, 5, 6, 7, 8:
			for i := 0; i < 2; i++ {
				mode := (program[pc] / (10 * int(math.Pow(10, float64(i+1))))) % 10
				if mode == 0 {
					// position mode
					params[i] = program[program[pc+i+1]]
				} else {
					// immediate mode
					params[i] = program[pc+i+1]
				}
			}
		}

		switch opcode {

		case 1:
			program[program[pc+3]] = params[0] + params[1]
			pc += 4
		case 2:
			program[program[pc+3]] = params[0] * params[1]
			pc += 4
		case 3:
			program[program[pc+1]] = inputValue
			pc += 2
		case 4:
			fmt.Printf("output: %d\n", program[program[pc+1]])
			pc += 2
		case 5:
			if params[0] != 0 {
				pc = params[1]
			} else {
				pc += 3
			}
		case 6:
			if params[0] == 0 {
				pc = params[1]
			} else {
				pc += 3
			}
		case 7:
			toStore := 0
			if params[0] < params[1] {
				toStore = 1
			}
			program[program[pc+3]] = toStore
			pc += 4
		case 8:
			toStore := 0
			if params[0] == params[1] {
				toStore = 1
			}
			program[program[pc+3]] = toStore
			pc += 4
		case 99:
			return
		default:
			log.Fatalf("unknown opcode pc(%d) program(%v)", pc, program)
		}
	}
}

func main() {

	input := "3,225,1,225,6,6,1100,1,238,225,104,0,1101,72,36,225,1101,87,26,225,2,144,13,224,101,-1872,224,224,4,224,102,8,223,223,1001,224,2,224,1,223,224,223,1102,66,61,225,1102,25,49,224,101,-1225,224,224,4,224,1002,223,8,223,1001,224,5,224,1,223,224,223,1101,35,77,224,101,-112,224,224,4,224,102,8,223,223,1001,224,2,224,1,223,224,223,1002,195,30,224,1001,224,-2550,224,4,224,1002,223,8,223,1001,224,1,224,1,224,223,223,1102,30,44,225,1102,24,21,225,1,170,117,224,101,-46,224,224,4,224,1002,223,8,223,101,5,224,224,1,224,223,223,1102,63,26,225,102,74,114,224,1001,224,-3256,224,4,224,102,8,223,223,1001,224,3,224,1,224,223,223,1101,58,22,225,101,13,17,224,101,-100,224,224,4,224,1002,223,8,223,101,6,224,224,1,224,223,223,1101,85,18,225,1001,44,7,224,101,-68,224,224,4,224,102,8,223,223,1001,224,5,224,1,223,224,223,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,7,677,226,224,102,2,223,223,1005,224,329,101,1,223,223,8,677,226,224,1002,223,2,223,1005,224,344,1001,223,1,223,1107,677,677,224,102,2,223,223,1005,224,359,1001,223,1,223,1107,226,677,224,102,2,223,223,1005,224,374,101,1,223,223,7,226,677,224,102,2,223,223,1005,224,389,101,1,223,223,8,226,677,224,1002,223,2,223,1005,224,404,101,1,223,223,1008,226,677,224,1002,223,2,223,1005,224,419,1001,223,1,223,107,677,677,224,102,2,223,223,1005,224,434,101,1,223,223,1108,677,226,224,1002,223,2,223,1006,224,449,101,1,223,223,1108,677,677,224,102,2,223,223,1006,224,464,101,1,223,223,1007,677,226,224,102,2,223,223,1006,224,479,101,1,223,223,1008,226,226,224,102,2,223,223,1006,224,494,101,1,223,223,108,226,226,224,1002,223,2,223,1006,224,509,101,1,223,223,107,226,226,224,102,2,223,223,1006,224,524,101,1,223,223,1107,677,226,224,102,2,223,223,1005,224,539,1001,223,1,223,108,226,677,224,1002,223,2,223,1005,224,554,101,1,223,223,1007,226,226,224,102,2,223,223,1005,224,569,101,1,223,223,8,226,226,224,102,2,223,223,1006,224,584,101,1,223,223,1008,677,677,224,1002,223,2,223,1005,224,599,1001,223,1,223,107,226,677,224,1002,223,2,223,1005,224,614,1001,223,1,223,1108,226,677,224,102,2,223,223,1006,224,629,101,1,223,223,7,677,677,224,1002,223,2,223,1005,224,644,1001,223,1,223,108,677,677,224,102,2,223,223,1005,224,659,101,1,223,223,1007,677,677,224,102,2,223,223,1006,224,674,101,1,223,223,4,223,99,226"

	var program []int
	for _, value := range strings.Split(input, ",") {
		asInt, _ := strconv.Atoi(value)
		program = append(program, asInt)
	}

	var inputValue int
	fmt.Print("Input Value: ")
	fmt.Scanf("%d", &inputValue)

	runProgram(program, inputValue)
}
