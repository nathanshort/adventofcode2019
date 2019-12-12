package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
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

func newIntcodeComputer(instructions string) intcodeComputer {
	c := intcodeComputer{}
	c.program = make(map[int64]int64)
	for index, value := range strings.Split(instructions, ",") {
		asInt, _ := strconv.ParseInt(value, 10, 64)
		c.program[int64(index)] = asInt
	}
	return c
}

func (c *intcodeComputer) run(input chan int64, output chan int64, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}
	defer close(output)

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

type point struct {
	x, y int64
}

const (
	up    = iota
	right = iota
	down  = iota
	left  = iota
)

const (
	black = 0
	white = 1
)

type grid struct {
	location  point
	direction int
	seen      map[point]int64
	colors    map[point]int64
}

func newGrid(origin point, direction int, originColor int64) grid {
	g := grid{location: origin, direction: direction}
	g.seen = make(map[point]int64)
	g.colors = make(map[point]int64)
	g.seen[origin] = 1
	g.colors[origin] = originColor
	return g
}

func (g *grid) setColorCurrent(color int64) {
	g.colors[g.location] = color
}

func (g grid) getColorCurrent() int64 {
	return g.colors[g.location]
}

func (g grid) render(writer io.Writer) {
	var minx int64 = math.MaxInt64
	var miny int64 = math.MaxInt64
	var maxx int64 = math.MinInt64
	var maxy int64 = math.MinInt64

	for pp := range g.seen {
		if pp.x < minx {
			minx = pp.x
		}
		if pp.y < miny {
			miny = pp.y
		}
		if pp.x > maxx {
			maxx = pp.x
		}
		if pp.y > maxy {
			maxy = pp.y
		}
	}

	image := image.NewRGBA(image.Rect(int(minx)-10, int(miny)-10, int(maxx)+10, int(maxy)+10))
	for pp := range g.seen {
		switch g.colors[pp] {
		case black:
			image.Set(int(pp.x), -int(pp.y), color.RGBA{0, 0, 0, 255})
		case white:
			image.Set(int(pp.x), -int(pp.y), color.RGBA{255, 255, 255, 255})
		}
	}
	png.Encode(writer, image)
}

func (g *grid) move(spaces int64) {
	switch g.direction {
	case up:
		g.location = point{x: g.location.x, y: g.location.y + spaces}
	case left:
		g.location = point{x: g.location.x - spaces, y: g.location.y}
	case down:
		g.location = point{x: g.location.x, y: g.location.y - spaces}
	case right:
		g.location = point{x: g.location.x + spaces, y: g.location.y}
	default:
		log.Fatalf("unknown direction")
	}

	g.seen[g.location] = g.seen[g.location] + 1
}

func (g *grid) turn(direction int) {

	var turnLeft = map[int]int{up: left, left: down, down: right, right: up}
	var turnRight = map[int]int{up: right, right: down, down: left, left: up}

	switch direction {
	case left:
		g.direction = turnLeft[g.direction]
	case right:
		g.direction = turnRight[g.direction]
	default:
		log.Fatalf("dont know how to turn %d", direction)
	}
}

func runIteration(computer intcodeComputer, panels *grid) {

	inputChan := make(chan int64, 1)
	outputChan := make(chan int64)
	inputChan <- panels.getColorCurrent()

	go computer.run(inputChan, outputChan, nil)

	outputInstructions := make([]int64, 2)
	instructionsSeen := 0
	for i := range outputChan {
		outputInstructions[instructionsSeen] = i
		instructionsSeen++
		if instructionsSeen == 2 {
			panels.setColorCurrent(outputInstructions[0])
			switch outputInstructions[1] {
			case 0:
				panels.turn(left)
			case 1:
				panels.turn(right)
			default:
				log.Fatalf("unknown direction")
			}
			panels.move(1)
			instructionsSeen = 0
			inputChan <- panels.getColorCurrent()
		}
	}
}

func main() {

	input := "3,8,1005,8,310,1106,0,11,0,0,0,104,1,104,0,3,8,102,-1,8,10,1001,10,1,10,4,10,108,1,8,10,4,10,1002,8,1,28,1,105,11,10,3,8,102,-1,8,10,1001,10,1,10,4,10,1008,8,0,10,4,10,102,1,8,55,3,8,102,-1,8,10,1001,10,1,10,4,10,108,0,8,10,4,10,1001,8,0,76,3,8,1002,8,-1,10,101,1,10,10,4,10,108,0,8,10,4,10,102,1,8,98,1,1004,7,10,1006,0,60,3,8,102,-1,8,10,1001,10,1,10,4,10,108,0,8,10,4,10,1002,8,1,127,2,1102,4,10,1,1108,7,10,2,1102,4,10,2,101,18,10,3,8,1002,8,-1,10,1001,10,1,10,4,10,1008,8,0,10,4,10,102,1,8,166,1006,0,28,3,8,1002,8,-1,10,101,1,10,10,4,10,108,1,8,10,4,10,101,0,8,190,1006,0,91,1,1108,5,10,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,1,10,4,10,1002,8,1,220,1,1009,14,10,2,1103,19,10,2,1102,9,10,2,1007,4,10,3,8,1002,8,-1,10,101,1,10,10,4,10,1008,8,1,10,4,10,101,0,8,258,2,3,0,10,1006,0,4,3,8,102,-1,8,10,1001,10,1,10,4,10,108,1,8,10,4,10,1001,8,0,286,1006,0,82,101,1,9,9,1007,9,1057,10,1005,10,15,99,109,632,104,0,104,1,21102,1,838479487636,1,21102,327,1,0,1106,0,431,21102,1,932813579156,1,21102,1,338,0,1106,0,431,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,21101,0,179318033447,1,21101,385,0,0,1105,1,431,21101,248037678275,0,1,21101,0,396,0,1105,1,431,3,10,104,0,104,0,3,10,104,0,104,0,21101,0,709496558348,1,21102,419,1,0,1105,1,431,21101,825544561408,0,1,21101,0,430,0,1106,0,431,99,109,2,22101,0,-1,1,21101,40,0,2,21102,462,1,3,21101,0,452,0,1106,0,495,109,-2,2105,1,0,0,1,0,0,1,109,2,3,10,204,-1,1001,457,458,473,4,0,1001,457,1,457,108,4,457,10,1006,10,489,1101,0,0,457,109,-2,2106,0,0,0,109,4,2101,0,-1,494,1207,-3,0,10,1006,10,512,21101,0,0,-3,22101,0,-3,1,22101,0,-2,2,21101,1,0,3,21102,531,1,0,1105,1,536,109,-4,2105,1,0,109,5,1207,-3,1,10,1006,10,559,2207,-4,-2,10,1006,10,559,22101,0,-4,-4,1106,0,627,21202,-4,1,1,21201,-3,-1,2,21202,-2,2,3,21102,578,1,0,1105,1,536,22101,0,1,-4,21101,1,0,-1,2207,-4,-2,10,1006,10,597,21102,0,1,-1,22202,-2,-1,-2,2107,0,-3,10,1006,10,619,21201,-1,0,1,21102,1,619,0,105,1,494,21202,-2,-1,-2,22201,-4,-2,-4,109,-5,2106,0,0"
	computer := newIntcodeComputer(input)

	{
		panels := newGrid(point{}, up, black)
		runIteration(computer, &panels)
		fmt.Printf("part 1: %d\n", len(panels.colors))
	}

	{
		panels := newGrid(point{}, up, white)
		runIteration(computer, &panels)
		f, _ := os.Create("/tmp/10.image")
		defer f.Close()
		panels.render(f)
	}

}
