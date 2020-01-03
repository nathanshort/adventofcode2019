package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"unicode"
	"unsafe"
)

type point struct {
	x, y int
}

type grid struct {
	points map[point]rune
	keys   map[point]rune
	doors  map[point]rune
	start  point
}

func newGrid() *grid {
	g := &grid{}
	g.points = make(map[point]rune)
	g.keys = make(map[point]rune)
	g.doors = make(map[point]rune)
	return g
}

func newGridFromScanner(scanner *bufio.Scanner) *grid {
	grid := newGrid()
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		x := 0
		for _, char := range line {
			point := point{x, y}
			grid.points[point] = char
			if unicode.IsUpper(char) {
				grid.doors[point] = char
			} else if unicode.IsLower(char) {
				grid.keys[point] = char
			} else if char == '@' {
				grid.start = point
				grid.points[point] = '.'
			}
			x++
		}
		y--
	}
	return grid
}

type pointAndUnlocked struct {
	point

	// this is a bitvector Z is bit 0; A is bit 25
	unlocked int
}

/// can pass a door or a key, as we will normalize to upper
func getBit(door rune) uint {
	return uint('Z' - unicode.ToUpper(door))
}

func (p pointAndUnlocked) isUnlocked(door rune) bool {
	shifted := 1 << getBit(door)
	return (p.unlocked & (shifted)) == shifted
}

func (p *pointAndUnlocked) unlock(door rune) {
	p.unlocked |= (1 << getBit(door))
}

func (p pointAndUnlocked) getSetBitCount() int {
	numSet := 0
	for i := uint(0); i < uint(unsafe.Sizeof(p.unlocked))*8; i++ {
		shifted := 1 << i
		if p.unlocked&(shifted) == shifted {
			numSet++
		}
	}
	return numSet
}

type pointAndUnlockedAndDistance struct {
	pointAndUnlocked
	distance int
}

func bfs(grid *grid) (int, bool) {

	visited := make(map[pointAndUnlocked]bool)
	start := pointAndUnlocked{point: grid.start}
	visited[start] = true

	queue := []pointAndUnlockedAndDistance{}
	queue = append(queue, pointAndUnlockedAndDistance{pointAndUnlocked: start, distance: 0})
	minDistance := math.MaxInt32

	for len(queue) != 0 {

		current := queue[0]
		queue = queue[1:]

		if current.getSetBitCount() == len(grid.keys) {
			minDistance = current.distance
			break
		}

		var movements = []point{{x: 0, y: 1}, {x: 0, y: -1}, {x: -1, y: 0}, {x: 1, y: 0}}
		for _, move := range movements {

			nextItem := pointAndUnlocked{point: point{x: current.point.x + move.x, y: current.point.y + move.y},
				unlocked: current.unlocked}
			if visited[nextItem] {
				continue
			}

			canMove := false
			nextPointValue, ok := grid.points[nextItem.point]
			if ok != true {
				// do nothing
			} else if nextPointValue == '.' {
				canMove = true
			} else if unicode.IsUpper(nextPointValue) {

				/// if the door is unlocked, then we can continue.
				/// if the door's key is not in this grid, then we can continue - under the understanding
				/// that we'll just wait and some other grid will open it.  this only comes into play
				/// in part2.  not sure if this is always legit...but it works
				/// map of keys would be quicker
				///
				if current.isUnlocked(nextPointValue) {
					canMove = true
				} else {
					keyIsInThisGrid := false
					keyValue := unicode.ToLower(nextPointValue)
					for _, k := range grid.keys {
						if k == keyValue {
							keyIsInThisGrid = true
							break
						}
					}
					canMove = !keyIsInThisGrid
				}
			} else if unicode.IsLower(nextPointValue) {
				canMove = true
				nextItem.unlock(nextPointValue)
			}

			if canMove {
				visited[nextItem] = true
				queue = append(queue, pointAndUnlockedAndDistance{pointAndUnlocked: nextItem, distance: current.distance + 1})
			}
		}
	}

	if minDistance != math.MaxInt32 {
		return minDistance, true
	}
	return 0, false
}

func part1() {

	grid := newGridFromScanner(bufio.NewScanner(os.Stdin))
	distance, ok := bfs(grid)
	if !ok {
		log.Fatalf("failed to find solution\n")
	} else {
		fmt.Printf("part 1 min distance %d\n", distance)
	}
}

func part2() {

	allDistances := 0
	for i := 1; i <= 4; i++ {
		file, _ := os.Open(fmt.Sprintf("input.%d", i))
		grid := newGridFromScanner(bufio.NewScanner(file))
		distance, ok := bfs(grid)
		if !ok {
			log.Fatalf("failed to find solution\n")
		} else {
			allDistances += distance
		}
	}
	fmt.Printf("part 2:%d\n", allDistances)
}

func main() {

	part1()
	part2()
}
