package main

import (
	"fmt"
	"strconv"
)

func part1(min int, max int) int {

	validPasswords := 0

	for i := min; i <= max; i++ {
		asString := strconv.Itoa(i)
		var lastSeenChar rune = '0'
		haveSeenDouble := false
		alwaysIncreasing := true

		for _, char := range asString {
			if char < lastSeenChar {
				alwaysIncreasing = false
				break
			} else if char == lastSeenChar {
				haveSeenDouble = true
			}
			lastSeenChar = char
		}
		if alwaysIncreasing && haveSeenDouble {
			validPasswords++
		}
	}

	return validPasswords
}

func part2(min int, max int) int {

	validPasswords := 0

	for i := min; i <= max; i++ {
		asString := strconv.Itoa(i)
		var lastSeenChar rune = '0'
		haveSeenDouble := false
		alwaysIncreasing := true
		repeatCount := 0
		validRepeat := false

		for index, char := range asString {
			if char < lastSeenChar {
				alwaysIncreasing = false
				break
			} else if char == lastSeenChar {
				haveSeenDouble = true
				repeatCount++
			}

			if repeatCount == 1 && (char != lastSeenChar || index == len(asString)-1) {
				validRepeat = true
			}
			if char != lastSeenChar {
				repeatCount = 0
			}

			lastSeenChar = char
		}
		if alwaysIncreasing && haveSeenDouble && validRepeat {
			validPasswords++
		}
	}
	return validPasswords

}

func main() {

	min := 265275
	max := 781584

	fmt.Printf("part 1: %d\n", part1(min, max))
	fmt.Printf("part 2: %d\n", part2(min, max))

}
