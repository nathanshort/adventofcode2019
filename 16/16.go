package main

import (
	"fmt"
	"math"
	"strconv"
)

/// clunky.  skips over signal elements in batches, depending on the pattern
/// if, if the pattern is 0 0 0 1 1 1 ..., then we'll skip 3 signals at a time
/// as they would be 0 multiplied anyway
/// prob a premature optimization
func phase(signal []int, pattern []int) []int {

	output := make([]int, len(signal))
	for iteration := 0; iteration < len(signal); iteration++ {

		patternPointer := 0
		iterationsThisPatternElement := 0
		haveShifted := false

		for i := 0; i < len(signal); {
			if !haveShifted {
				haveShifted = true
				if iteration == 0 {
					patternPointer++
					continue
				} else {
					iterationsThisPatternElement++
				}
			}
			if pattern[patternPointer] == 0 {
				patternPointer = (patternPointer + 1) % len(pattern)
				i += (iteration + 1 - iterationsThisPatternElement)
				iterationsThisPatternElement = 0
			} else {
				output[iteration] += signal[i] * pattern[patternPointer]
				i++
				if iterationsThisPatternElement >= iteration {
					iterationsThisPatternElement = 0
					patternPointer = (patternPointer + 1) % len(pattern)
				} else {
					iterationsThisPatternElement++
				}
			}
		}
		output[iteration] = int(math.Abs(float64(output[iteration] % 10)))
	}
	return output
}

func part1(signal []int, pattern []int) {
	for i := 0; i < 100; i++ {
		signal = phase(signal, pattern)
	}
	fmt.Printf("part 1: %v\n", signal[0:8])
}

/// take advantage of a couple of things
/// 1) as we are shifted over so far, the pattern will end up being all 1's - so we can just run additions
/// 2) the value at i is i + sum( i+1..end )
func phase2(signal []int, offset int) []int {
	output := make([]int, len(signal))
	sum := 0

	for i := len(signal) - 1; i >= offset; i-- {
		sum += signal[i]
		output[i] = sum % 10
	}
	return output
}

func part2(signal []int, pattern []int, offset int) {

	fullsignal := []int{}
	for i := 0; i < 10000; i++ {
		fullsignal = append(fullsignal, signal...)
	}
	for i := 0; i < 100; i++ {
		fullsignal = phase2(fullsignal, offset)
	}
	fmt.Printf("part 2: %v\n", fullsignal[offset:offset+8])
}

func main() {

	input := "59708372326282850478374632294363143285591907230244898069506559289353324363446827480040836943068215774680673708005813752468017892971245448103168634442773462686566173338029941559688604621181240586891859988614902179556407022792661948523370366667688937217081165148397649462617248164167011250975576380324668693910824497627133242485090976104918375531998433324622853428842410855024093891994449937031688743195134239353469076295752542683739823044981442437538627404276327027998857400463920633633578266795454389967583600019852126383407785643022367809199144154166725123539386550399024919155708875622641704428963905767166129198009532884347151391845112189952083025"

	signal := []int{}
	for _, v := range input {
		signal = append(signal, int(v-48))
	}
	pattern := []int{0, 1, 0, -1}
	offset, _ := strconv.Atoi(input[:7])

	part1(signal, pattern)
	part2(signal, pattern, offset)
}
