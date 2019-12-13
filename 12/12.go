package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type moon struct {
	x, y, z, xv, yv, zv int
}

func (m *moon) applyVelocity() {
	m.x += m.xv
	m.y += m.yv
	m.z += m.zv
}

func (m moon) potential() int {
	return int(math.Abs(float64(m.x)) + math.Abs(float64(m.y)) + math.Abs(float64(m.z)))
}

func (m moon) kinetic() int {
	return int(math.Abs(float64(m.xv)) + math.Abs(float64(m.yv)) + math.Abs(float64(m.zv)))
}

func (m moon) totalEnergy() int {
	return m.potential() * m.kinetic()
}

func applyGravity(val1 int, val2 int, target1 *int, target2 *int) {
	if val1 == val2 {
		return
	} else if val1 < val2 {
		*target1++
		*target2--
	} else {
		*target1--
		*target2++
	}
}

func greatestCommonDivisor(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

/// stolen from the internet
func leastCommonMultiple(a, b int64, integers ...int64) int64 {
	result := a * b / greatestCommonDivisor(a, b)
	for i := 0; i < len(integers); i++ {
		result = leastCommonMultiple(result, integers[i])
	}
	return result
}

func main() {
	var moons []*moon
	var initialMoons []*moon
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var x, y, z int
		fmt.Sscanf(scanner.Text(), "<x=%d, y=%d, z=%d>", &x, &y, &z)
		moons = append(moons, &moon{x: x, y: y, z: z})
		initialMoons = append(initialMoons, &moon{x: x, y: y, z: z})
	}

	/// number of iterations before velocity for all moons, at the specified axis, is 0
	var xv0, yv0, zv0 int

	numIterations := 1000000
	for iteration := 0; iteration < numIterations; iteration++ {
		for i := 0; i < len(moons); i++ {
			for j := i + 1; j < len(moons); j++ {
				m1 := moons[i]
				m2 := moons[j]
				applyGravity(m1.x, m2.x, &m1.xv, &m2.xv)
				applyGravity(m1.y, m2.y, &m1.yv, &m2.yv)
				applyGravity(m1.z, m2.z, &m1.zv, &m2.zv)
			}
		}
		for _, m := range moons {
			m.applyVelocity()
		}

		/// find if all moons are terminal in velocity at any axis
		foundxv0 := true
		foundyv0 := true
		foundzv0 := true
		for _, m := range moons {
			if foundxv0 && m.xv != 0 {
				foundxv0 = false
			}
			if foundyv0 && m.yv != 0 {
				foundyv0 = false
			}
			if foundzv0 && m.zv != 0 {
				foundzv0 = false
			}
		}
		if xv0 == 0 && foundxv0 {
			xv0 = iteration + 1
		}
		if yv0 == 0 && foundyv0 {
			yv0 = iteration + 1
		}
		if zv0 == 0 && foundzv0 {
			zv0 = iteration + 1
		}
	}

	/// change numIterations to requested numIterations prior
	/// to running part 1
	if false {
		totalEnergy := 0
		for _, m := range moons {
			totalEnergy += m.totalEnergy()
		}
		fmt.Printf("part 1: %d\n", totalEnergy)
	}

	/// part 2
	if true {

		if xv0 == 0 || yv0 == 0 || zv0 == 0 {
			log.Fatalf("%d is not enough iterations to find cycle\n", numIterations)
		}

		/// *v0 is where all moons stopped moving on that axis. 2 x lcm( *v0 )
		/// is when they have returned back to their starting point
		fmt.Printf("part 2: %v\n", 2*leastCommonMultiple(int64(xv0), int64(yv0), int64(zv0)))
	}
}
