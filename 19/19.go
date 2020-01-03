package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type intcodeComputer struct {
	pc           int64
	relativeBase int64
	program      map[int64]int64
}

func (c intcodeComputer) mode(pcOffset int64) int64 {
	return c.program[c.pc] / (10 * int64(math.Pow(10, float64(pcOffset)))) % 10
}

func (c intcodeComputer) read(pcOffset int64) int64 {
	switch c.mode(pcOffset) {
	case 0:
		return c.program[c.program[c.pc+pcOffset]]
	case 1:
		return c.program[c.pc+pcOffset]
	case 2:
		return c.program[c.program[c.pc+pcOffset]+c.relativeBase]
	default:
		log.Fatalf("bad read mode")
		return -1
	}
}

func (c *intcodeComputer) write(pcOffset int64, value int64) {
	switch c.mode(pcOffset) {
	case 0:
		c.program[c.program[c.pc+pcOffset]] = value
	case 2:
		c.program[c.program[c.pc+pcOffset]+c.relativeBase] = value
	default:
		log.Fatalf("bad write mode")
	}
}

func newIntcodeComputer(instructions string) *intcodeComputer {
	c := &intcodeComputer{}
	c.program = make(map[int64]int64)
	for index, value := range strings.Split(instructions, ",") {
		c.program[int64(index)], _ = strconv.ParseInt(value, 10, 64)
	}
	return c
}

func (c *intcodeComputer) run(input chan int64, output chan int64) {

	for {
		opcode := c.program[c.pc]%10 + c.program[c.pc]/10%10*10
		switch opcode {
		case 1:
			c.write(3, c.read(1)+c.read(2))
			c.pc += 4
		case 2:
			c.write(3, c.read(1)*c.read(2))
			c.pc += 4
		case 3:
			c.write(1, <-input)
			c.pc += 2
		case 4:
			output <- c.read(1)
			c.pc += 2
		case 5:
			if c.read(1) != 0 {
				c.pc = c.read(2)
			} else {
				c.pc += 3
			}
		case 6:
			if c.read(1) == 0 {
				c.pc = c.read(2)
			} else {
				c.pc += 3
			}
		case 7:
			toStore := int64(0)
			if c.read(1) < c.read(2) {
				toStore = 1
			}
			c.write(3, toStore)
			c.pc += 4
		case 8:
			toStore := int64(0)
			if c.read(1) == c.read(2) {
				toStore = 1
			}
			c.write(3, toStore)
			c.pc += 4
		case 9:
			c.relativeBase += c.read(1)
			c.pc += 2
		case 99:
			return
		default:
			log.Fatalf("unknown opcode pc(%d) program(%v)", c.pc, c.program)
		}
	}
}

func part1(input string) {
	numPulled := int64(0)
	for y := int64(0); y < 50; y++ {
		for x := int64(0); x < 50; x++ {
			inputChan := make(chan int64)
			outputChan := make(chan int64)
			computer := newIntcodeComputer(input)
			go computer.run(inputChan, outputChan)
			inputChan <- x
			inputChan <- y
			result := <-outputChan
			numPulled += result
		}
	}
	fmt.Printf("part 1: %d\n", numPulled)
}

type point struct {
	x, y int64
}

type grid struct {
	points map[point]bool
}

func newGrid() *grid {
	g := &grid{}
	g.points = make(map[point]bool)
	return g
}

func part2(input string) {

	grid := newGrid()

	/// we dont always need the x to start at 0.  start at the x of the first hit
	/// last time through
	xStart := int64(0)

	/// and we can jump a bunch of x by the width of the beam last y. no need
	/// to check the points in between
	xWidth := int64(0)

	for y := int64(0); ; y++ {
		xStartThisY := int64(-1)
		hitsThisY := 0

		/// 5000 is ... somewhat arbitrary. any way to figure out what it should be?
		for x := xStart; x < 5000; x++ {
			computer := newIntcodeComputer(input)
			inputChan := make(chan int64)
			outputChan := make(chan int64)
			go computer.run(inputChan, outputChan)
			inputChan <- x
			inputChan <- y

			result := <-outputChan
			if result == 1 {
				grid.points[point{x, y}] = true
				hitsThisY++

				/// first hit this y
				if xStartThisY == -1 {
					xStartThisY = x
					xStart = x
					x += xWidth - 1
				}
			} else if hitsThisY != 0 {
				/// we have already had hits this y, but now we are not
				/// receiving a hit.  thus, we can end this round
				xWidth = x - xStartThisY
				break
			}
		}

		if grid.points[point{xStartThisY + 99, y - 99}] {
			answer := xStartThisY*10000 + (y - 99)
			fmt.Printf("part 2: %d\n", answer)
			return
		}
	}
}

func main() {

	input := "109,424,203,1,21101,11,0,0,1106,0,282,21102,18,1,0,1106,0,259,2101,0,1,221,203,1,21101,31,0,0,1106,0,282,21102,1,38,0,1106,0,259,21002,23,1,2,22102,1,1,3,21102,1,1,1,21101,57,0,0,1106,0,303,2101,0,1,222,21002,221,1,3,21001,221,0,2,21102,259,1,1,21102,80,1,0,1106,0,225,21102,1,79,2,21101,0,91,0,1106,0,303,2102,1,1,223,21001,222,0,4,21102,259,1,3,21101,225,0,2,21102,1,225,1,21101,0,118,0,1105,1,225,21002,222,1,3,21101,118,0,2,21101,0,133,0,1106,0,303,21202,1,-1,1,22001,223,1,1,21102,1,148,0,1105,1,259,1202,1,1,223,20102,1,221,4,20101,0,222,3,21102,1,22,2,1001,132,-2,224,1002,224,2,224,1001,224,3,224,1002,132,-1,132,1,224,132,224,21001,224,1,1,21102,1,195,0,105,1,109,20207,1,223,2,21002,23,1,1,21101,-1,0,3,21102,214,1,0,1106,0,303,22101,1,1,1,204,1,99,0,0,0,0,109,5,2101,0,-4,249,22101,0,-3,1,22102,1,-2,2,21201,-1,0,3,21101,0,250,0,1105,1,225,22101,0,1,-4,109,-5,2105,1,0,109,3,22107,0,-2,-1,21202,-1,2,-1,21201,-1,-1,-1,22202,-1,-2,-2,109,-3,2106,0,0,109,3,21207,-2,0,-1,1206,-1,294,104,0,99,22102,1,-2,-2,109,-3,2106,0,0,109,5,22207,-3,-4,-1,1206,-1,346,22201,-4,-3,-4,21202,-3,-1,-1,22201,-4,-1,2,21202,2,-1,-1,22201,-4,-1,1,22102,1,-2,3,21102,343,1,0,1106,0,303,1105,1,415,22207,-2,-3,-1,1206,-1,387,22201,-3,-2,-3,21202,-2,-1,-1,22201,-3,-1,3,21202,3,-1,-1,22201,-3,-1,2,21201,-4,0,1,21102,384,1,0,1105,1,303,1106,0,415,21202,-4,-1,-4,22201,-4,-3,-4,22202,-3,-2,-2,22202,-2,-4,-4,22202,-3,-2,-3,21202,-4,-1,-2,22201,-3,-2,1,22101,0,1,-4,109,-5,2106,0,0"
	part1(input)
	part2(input)

}
