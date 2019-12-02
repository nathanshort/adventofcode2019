package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func fuel(mass int) int {
	return int(math.Floor(float64(mass/3))) - 2
}

func part1(masses []int) {
	var sum int = 0
	for _, mass := range masses {
		sum += fuel(mass)
	}
	fmt.Printf("part1 fuel is %d\n", sum)
}

func part2(masses []int) {
	var sum int = 0

	for _, mass := range masses {
		current := mass

		for {
			fuelNeeded := fuel(current)
			if fuelNeeded > 0 {
				sum += fuelNeeded
				current = fuelNeeded
			} else {
				break
			}
		}
	}
	fmt.Printf("part2 fuel is %d\n", sum)

}

func main() {
	var masses []int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		mass, _ := strconv.Atoi(scanner.Text())
		masses = append(masses, mass)
	}

	part1(masses)
	part2(masses)
}
