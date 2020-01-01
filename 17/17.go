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

type robot struct {
	where       point
	orientation int64
}

type grid struct {
	points map[point]int64
	robot
}

func newGrid() *grid {
	g := &grid{}
	g.points = make(map[point]int64)
	return g
}

func (g grid) print() {
	for y := 10; y >= -50; y-- {
		for x := -20; x <= 60; x++ {
			atMove, ok := g.points[point{x: x, y: y}]
			toPrint := "."
			if ok == true {
				switch atMove {
				case scaffold:
					toPrint = "#"
				}
			}
			fmt.Printf("%s", toPrint)
		}
		fmt.Printf("\n")
	}
}

type point struct {
	x, y int
}

const (
	scaffold = 35
	open     = 46
	newline  = 10
	up       = 94
	right    = 62
	left     = 60
	down     = 118
)

func part1(input string) *grid {
	computer := newIntcodeComputer(input)
	inputChan := make(chan int64)
	outputChan := make(chan int64)

	go computer.run(inputChan, outputChan)

	grid := newGrid()
	currentPoint := point{}

	for i := range outputChan {
		dx, dy := 0, 0
		switch i {
		case open:
			dx = 1
		case newline:
			dy, dx = -1, -currentPoint.x
		case up, down, left, right:
			grid.where = currentPoint
			grid.orientation = i
			dx = 1
		case scaffold:
			grid.points[currentPoint] = i
			dx = 1
		default:
			log.Fatalf("unknown %d\n", i)
		}
		currentPoint = point{x: currentPoint.x + dx, y: currentPoint.y + dy}
	}

	alignments := 0
	for thePoint, v := range grid.points {
		if v != scaffold {
			continue
		}
		intersection := true
		moves := []point{{1, 0}, {0, -1}, {-1, 0}, {0, 1}}
		for _, move := range moves {
			atMove, ok := grid.points[point{x: thePoint.x + move.x, y: thePoint.y + move.y}]
			if ok == false || atMove != scaffold {
				intersection = false
				break
			}
		}
		if intersection {
			alignments += int(math.Abs(float64(thePoint.x)) * math.Abs(float64(thePoint.y)))
		}
	}
	fmt.Printf("part 1: %d\n", alignments)
	return grid
}

func part2(input string, grid *grid) {

	type turnOption struct {
		turn, orientation int64
		dx, dy            int
	}

	/// this prob was made a lot easier as there is only one way to go at each turn.
	/// so, just figure out which way to turn, then, go that way for as long as possible
	for {
		var options []turnOption
		switch grid.orientation {
		case up:
			options = []turnOption{
				{turn: right, orientation: right, dx: 1, dy: 0},
				{turn: left, orientation: left, dx: -1, dy: 0}}
		case down:
			options = []turnOption{
				{turn: right, orientation: left, dx: -1, dy: 0},
				{turn: left, orientation: right, dx: 1, dy: 0}}
		case left:
			options = []turnOption{
				{turn: right, orientation: up, dx: 0, dy: 1},
				{turn: left, orientation: down, dx: 0, dy: -1}}
		case right:
			options = []turnOption{
				{turn: right, orientation: down, dx: 0, dy: -1},
				{turn: left, orientation: up, dx: 0, dy: 1}}
		}

		var option *turnOption
		for _, o := range options {
			nextPoint := point{x: grid.robot.where.x + o.dx, y: grid.robot.where.y + o.dy}
			if _, ok := grid.points[nextPoint]; ok != false {
				option = &o
				break
			}
		}
		if option == nil {
			break
		}

		/// we found a direction to move.  turn that direction, then, move as far as possible
		grid.orientation = option.orientation
		moves := 0
		for {
			nextPoint := point{x: grid.robot.where.x + option.dx, y: grid.robot.where.y + option.dy}
			if _, ok := grid.points[nextPoint]; ok == false {
				break
			}
			grid.where = nextPoint
			moves++
		}
		//		fmt.Printf("%d %d\n", option.turn, moves)
	}

	/// generated by looking at the move output ( currently commented out ) a couple lines above
	cmd := "A,A,B,C,A,C,A,B,C,B\nR,12,L,8,R,6\nR,12,L,6,R,6,R,8,R,6\nL,8,R,8,R,6,R,12\nn\n"

	computer := newIntcodeComputer(input)
	inputChan := make(chan int64, len(cmd))
	outputChan := make(chan int64)
	for _, c := range cmd {
		inputChan <- int64(c)
	}
	go computer.run(inputChan, outputChan)

	/// not quite sure whats up here.  we're getting grid output even when
	/// running in this mode.  program description didnt really mention anything about that.
	/// anyway - the last output is the value that we want
	var last int64 = 0
	for i := range outputChan {
		last = i
	}
	fmt.Printf("part 2: %d\n", last)
}

func main() {

	input := "1,330,331,332,109,3468,1102,1182,1,16,1101,0,1479,24,101,0,0,570,1006,570,36,1002,571,1,0,1001,570,-1,570,1001,24,1,24,1106,0,18,1008,571,0,571,1001,16,1,16,1008,16,1479,570,1006,570,14,21102,58,1,0,1105,1,786,1006,332,62,99,21101,0,333,1,21101,73,0,0,1105,1,579,1101,0,0,572,1102,1,0,573,3,574,101,1,573,573,1007,574,65,570,1005,570,151,107,67,574,570,1005,570,151,1001,574,-64,574,1002,574,-1,574,1001,572,1,572,1007,572,11,570,1006,570,165,101,1182,572,127,101,0,574,0,3,574,101,1,573,573,1008,574,10,570,1005,570,189,1008,574,44,570,1006,570,158,1105,1,81,21101,0,340,1,1106,0,177,21102,477,1,1,1106,0,177,21101,0,514,1,21102,176,1,0,1105,1,579,99,21101,0,184,0,1106,0,579,4,574,104,10,99,1007,573,22,570,1006,570,165,1001,572,0,1182,21102,375,1,1,21102,211,1,0,1105,1,579,21101,1182,11,1,21102,222,1,0,1106,0,979,21101,0,388,1,21101,233,0,0,1106,0,579,21101,1182,22,1,21102,244,1,0,1106,0,979,21101,0,401,1,21101,0,255,0,1105,1,579,21101,1182,33,1,21102,266,1,0,1105,1,979,21102,414,1,1,21102,1,277,0,1106,0,579,3,575,1008,575,89,570,1008,575,121,575,1,575,570,575,3,574,1008,574,10,570,1006,570,291,104,10,21101,1182,0,1,21102,313,1,0,1106,0,622,1005,575,327,1101,1,0,575,21101,0,327,0,1105,1,786,4,438,99,0,1,1,6,77,97,105,110,58,10,33,10,69,120,112,101,99,116,101,100,32,102,117,110,99,116,105,111,110,32,110,97,109,101,32,98,117,116,32,103,111,116,58,32,0,12,70,117,110,99,116,105,111,110,32,65,58,10,12,70,117,110,99,116,105,111,110,32,66,58,10,12,70,117,110,99,116,105,111,110,32,67,58,10,23,67,111,110,116,105,110,117,111,117,115,32,118,105,100,101,111,32,102,101,101,100,63,10,0,37,10,69,120,112,101,99,116,101,100,32,82,44,32,76,44,32,111,114,32,100,105,115,116,97,110,99,101,32,98,117,116,32,103,111,116,58,32,36,10,69,120,112,101,99,116,101,100,32,99,111,109,109,97,32,111,114,32,110,101,119,108,105,110,101,32,98,117,116,32,103,111,116,58,32,43,10,68,101,102,105,110,105,116,105,111,110,115,32,109,97,121,32,98,101,32,97,116,32,109,111,115,116,32,50,48,32,99,104,97,114,97,99,116,101,114,115,33,10,94,62,118,60,0,1,0,-1,-1,0,1,0,0,0,0,0,0,1,24,22,0,109,4,2102,1,-3,586,21002,0,1,-1,22101,1,-3,-3,21101,0,0,-2,2208,-2,-1,570,1005,570,617,2201,-3,-2,609,4,0,21201,-2,1,-2,1105,1,597,109,-4,2106,0,0,109,5,1202,-4,1,630,20101,0,0,-2,22101,1,-4,-4,21102,0,1,-3,2208,-3,-2,570,1005,570,781,2201,-4,-3,652,21002,0,1,-1,1208,-1,-4,570,1005,570,709,1208,-1,-5,570,1005,570,734,1207,-1,0,570,1005,570,759,1206,-1,774,1001,578,562,684,1,0,576,576,1001,578,566,692,1,0,577,577,21101,702,0,0,1105,1,786,21201,-1,-1,-1,1105,1,676,1001,578,1,578,1008,578,4,570,1006,570,724,1001,578,-4,578,21101,0,731,0,1105,1,786,1105,1,774,1001,578,-1,578,1008,578,-1,570,1006,570,749,1001,578,4,578,21101,756,0,0,1106,0,786,1105,1,774,21202,-1,-11,1,22101,1182,1,1,21102,1,774,0,1105,1,622,21201,-3,1,-3,1105,1,640,109,-5,2106,0,0,109,7,1005,575,802,20101,0,576,-6,20101,0,577,-5,1106,0,814,21102,0,1,-1,21101,0,0,-5,21102,0,1,-6,20208,-6,576,-2,208,-5,577,570,22002,570,-2,-2,21202,-5,51,-3,22201,-6,-3,-3,22101,1479,-3,-3,1202,-3,1,843,1005,0,863,21202,-2,42,-4,22101,46,-4,-4,1206,-2,924,21101,1,0,-1,1106,0,924,1205,-2,873,21102,1,35,-4,1106,0,924,1201,-3,0,878,1008,0,1,570,1006,570,916,1001,374,1,374,1201,-3,0,895,1101,0,2,0,1201,-3,0,902,1001,438,0,438,2202,-6,-5,570,1,570,374,570,1,570,438,438,1001,578,558,922,20101,0,0,-4,1006,575,959,204,-4,22101,1,-6,-6,1208,-6,51,570,1006,570,814,104,10,22101,1,-5,-5,1208,-5,39,570,1006,570,810,104,10,1206,-1,974,99,1206,-1,974,1102,1,1,575,21101,973,0,0,1105,1,786,99,109,-7,2106,0,0,109,6,21102,0,1,-4,21102,1,0,-3,203,-2,22101,1,-3,-3,21208,-2,82,-1,1205,-1,1030,21208,-2,76,-1,1205,-1,1037,21207,-2,48,-1,1205,-1,1124,22107,57,-2,-1,1205,-1,1124,21201,-2,-48,-2,1106,0,1041,21102,1,-4,-2,1106,0,1041,21102,-5,1,-2,21201,-4,1,-4,21207,-4,11,-1,1206,-1,1138,2201,-5,-4,1059,1201,-2,0,0,203,-2,22101,1,-3,-3,21207,-2,48,-1,1205,-1,1107,22107,57,-2,-1,1205,-1,1107,21201,-2,-48,-2,2201,-5,-4,1090,20102,10,0,-1,22201,-2,-1,-2,2201,-5,-4,1103,2102,1,-2,0,1105,1,1060,21208,-2,10,-1,1205,-1,1162,21208,-2,44,-1,1206,-1,1131,1105,1,989,21102,439,1,1,1105,1,1150,21102,477,1,1,1105,1,1150,21102,1,514,1,21101,1149,0,0,1105,1,579,99,21101,0,1157,0,1106,0,579,204,-2,104,10,99,21207,-3,22,-1,1206,-1,1138,1201,-5,0,1176,1202,-4,1,0,109,-6,2105,1,0,22,7,44,1,5,1,40,7,3,1,40,1,3,1,1,1,3,1,22,9,9,1,3,1,1,1,3,1,22,1,7,1,9,1,3,1,1,1,3,1,22,1,7,1,1,13,1,1,3,1,22,1,7,1,1,1,7,1,5,1,3,1,22,1,7,1,1,1,7,1,3,7,22,1,7,1,1,1,7,1,5,1,26,7,1,9,1,9,30,1,3,1,5,1,7,1,1,1,30,1,3,1,5,1,7,1,1,1,30,1,3,1,5,1,7,1,1,1,30,1,3,1,1,13,1,1,9,7,14,1,3,1,1,1,3,1,9,1,9,1,5,1,14,1,3,1,1,1,3,1,9,9,1,1,5,1,14,1,3,1,1,1,3,1,17,1,1,1,5,1,14,1,3,7,17,1,1,1,5,1,14,1,5,1,21,1,1,1,5,1,14,1,5,1,21,1,1,1,5,1,14,1,5,1,21,1,1,1,5,1,14,7,11,13,1,9,38,1,3,1,3,1,3,1,38,1,3,1,3,1,3,1,38,1,3,1,3,1,3,1,38,1,3,1,3,9,34,1,3,1,7,1,3,1,34,13,3,1,38,1,11,1,32,7,11,1,32,1,17,1,32,1,5,13,32,1,5,1,44,1,5,1,44,1,5,1,44,1,5,1,44,1,5,1,44,7,12"
	grid := part1(input)
	part2("2"+input[1:], grid)
}

/*

A = R,12,L,8,R,6\n
B = R,12,L,6,R,6,R,8,R,6\n
C = L,8,R,8,R,6,R,12,\n
CMD = A,A,B,C,A,C,A,B,C,B

62 12
60 8
62 6   A

62 12  A
60 8
62 6

62 12    B
60 6
62 6
62 8
62 6

60 8
62 8
62 6   C
62 12

62 12    A
60 8
62 6

60 8
62 8
62 6   C
62 12

62 12
60 8     A
62 6

62 12
60 6
62 6
62 8
62 6    B

60 8
62 8
62 6   C
62 12

62 12
60 6
62 6
62 8    B
62 6

*/
