package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

/// run the program, reading inputs from the input chan, and sending outputs to the output chan.
/// if the waitgroup is not null then decrement the waitgroup counter upon program termination
func runProgram(instructions []int, input chan int, output chan int, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

	program := make([]int, len(instructions))
	copy(program, instructions)
	pc := 0

	for {
		opcode := program[pc]%10 + program[pc]/10%10*10
		params := make([]int, 2)

		numOperands := 0
		switch opcode {
		case 1, 2, 5, 6, 7, 8:
			numOperands = 2
		case 4:
			numOperands = 1
		}

		for i := 0; i < numOperands; i++ {
			mode := (program[pc] / (10 * int(math.Pow(10, float64(i+1))))) % 10
			if mode == 0 {
				// position mode
				params[i] = program[program[pc+i+1]]
			} else {
				// immediate mode
				params[i] = program[pc+i+1]
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
			program[program[pc+1]] = <-input
			pc += 2
		case 4:
			output <- params[0]
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

/// find all 5 digit numbers with
/// 1) all unique digits
/// 2) no digits < minDigit
/// 3) no digits > maxDigit
func uniquePhaseSettings(minDigit int, maxDigit int) [][]int {
	var unique [][]int
	for attempt := 0; attempt < 100000; attempt++ {
		var digits []int
		digitsSeen := make(map[int]bool)
		allUnique := true
		for i := 4; i >= 0; i-- {
			digit := attempt / int(math.Pow(10, float64(i))) % 10
			if _, ok := digitsSeen[digit]; digit > maxDigit || digit < minDigit || ok != false {
				allUnique = false
				break
			}
			digitsSeen[digit] = true
			digits = append(digits, digit)
		}
		if allUnique {
			unique = append(unique, digits)
		}
	}
	return unique
}

func part1(program []int) {
	const numAmplifiers = 5
	maxOutput := 0
	for _, setting := range uniquePhaseSettings(0, 4) {
		lastOutput := 0
		for i := 0; i < numAmplifiers; i++ {
			input := make(chan int, 2)
			input <- setting[i]
			input <- lastOutput
			output := make(chan int, 1)
			runProgram(program, input, output, nil)
			lastOutput = <-output
		}
		if lastOutput > maxOutput {
			maxOutput = lastOutput
		}
	}
	fmt.Printf("part 1 max output: %d\n", maxOutput)
}

func part2(program []int) {

	largest := 0
	numAmplifiers := 5

	for _, setting := range uniquePhaseSettings(5, 9) {

		/// setup pipeline whereby the output of one amplifier
		/// is the input of the next amp.  feedback the last amp
		/// into the first
		var channels []chan int
		for i := 0; i < numAmplifiers; i++ {
			channel := make(chan int, 2)
			channel <- setting[i]
			channels = append(channels, channel)
		}
		/// seed amp A with the initial input
		channels[0] <- 0

		var wg sync.WaitGroup
		wg.Add(numAmplifiers)
		for i := 0; i < numAmplifiers; i++ {
			go runProgram(program, channels[i], channels[(i+1)%numAmplifiers], &wg)
		}
		wg.Wait()

		/// pull the result from the output of the last amp
		result := <-channels[0]
		if result > largest {
			largest = result
		}
	}

	fmt.Printf("part 2 max output: %d\n", largest)
}

func main() {

	input := "3,8,1001,8,10,8,105,1,0,0,21,38,63,76,93,118,199,280,361,442,99999,3,9,101,3,9,9,102,3,9,9,101,4,9,9,4,9,99,3,9,1002,9,2,9,101,5,9,9,1002,9,5,9,101,5,9,9,1002,9,4,9,4,9,99,3,9,101,2,9,9,102,3,9,9,4,9,99,3,9,101,2,9,9,102,5,9,9,1001,9,5,9,4,9,99,3,9,102,4,9,9,1001,9,3,9,1002,9,5,9,101,2,9,9,1002,9,2,9,4,9,99,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,2,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,99,3,9,101,1,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,99"

	var program []int
	for _, value := range strings.Split(input, ",") {
		asInt, _ := strconv.Atoi(value)
		program = append(program, asInt)
	}

	part1(program)
	part2(program)

}
