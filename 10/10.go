package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type point struct {
	x, y int
}

/// return distance between 2 points
func (p point) distance(other point) float64 {
	return math.Sqrt(math.Pow(float64(p.x-other.x), 2) + math.Pow(float64(p.y-other.y), 2))
}

/// used to keep slopes as integer fractions ( 1/3 ) instead of floats ( .3333 )
/// who knows if its really useful or if we could have just stored the float slope
type slopeAndDirection struct {
	dy, dx, xMag, yMag int
}

func (s slopeAndDirection) key() string {
	return fmt.Sprintf("%d/%d-%d,%d", s.dy, s.dx, s.xMag, s.yMag)
}

/// as we're storing slope as a fraction, we use this to normalize down to the gcd, such
/// that a slope of 2/5 is equal to one of 4/10
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func (p point) slopeAndDir(other point) slopeAndDirection {

	dy := p.y - other.y
	dx := p.x - other.x
	gcd := gcd(dy, dx)
	dy /= gcd
	dx /= gcd

	xMag := 0
	if p.x < other.x {
		xMag = -1
	} else if p.x > other.x {
		xMag = 1
	}

	yMag := 0
	if p.y < other.y {
		yMag = -1
	} else if p.y > other.y {
		yMag = 1
	}
	return slopeAndDirection{dy: dy, dx: dx, xMag: xMag, yMag: yMag}
}

/// returns list of vaporized points, in order of vaporization
func vaporize(bestPoint point, allAsteroids []point) []point {

	/// distance to bestPoint, keyed by other asteriod
	distances := make(map[point]float64)

	/// list of points at each angle
	angles := make(map[float64][]point)
	for _, a := range allAsteroids {
		if bestPoint == a {
			continue
		}
		/// change the coordinate system to treat bestPoint as 0,0
		/// scale existing points relative to the new origin
		scaledAsteroid := point{x: a.x - bestPoint.x, y: bestPoint.y - a.y}
		distances[scaledAsteroid] = bestPoint.distance(a)

		/// atan2 will give us an angle [ pi, -pi ] relative to the x axis, but what we really want is
		/// an angle such that vertical up is 0, and as we go clockwise one revolution gets us 2pi.
		/// so thats what this is doing.  taking the atan2 value and adjusting to get us what we want
		atan := math.Atan2(float64(scaledAsteroid.y), float64(scaledAsteroid.x))
		switch {
		case scaledAsteroid.x < 0 && scaledAsteroid.y >= 0:
			atan = 5*math.Pi/2 - atan
		default:
			atan = math.Pi/2 - atan
		}
		angles[atan] = append(angles[atan], scaledAsteroid)
	}

	var sortedAngles []float64
	for k := range angles {
		sortedAngles = append(sortedAngles, k)
	}
	sort.Slice(sortedAngles, func(i, j int) bool {
		return sortedAngles[i] < sortedAngles[j]
	})

	numVaporized := 0
	var pointsVaporized []point
	for i := 0; numVaporized < 200; i = (i + 1) % len(sortedAngles) {
		/// which points are at this angle
		eligiblePoints := angles[sortedAngles[i]]
		bestDistance := math.MaxFloat64
		var candidatePoint point
		/// which point has smallest distance to bestPoint
		for _, p := range eligiblePoints {
			if d, ook := distances[p]; ook == true {
				if d < bestDistance {
					bestDistance = distances[p]
					candidatePoint = p
				}
			}
		}
		if bestDistance != math.MaxFloat64 {
			numVaporized++
			/// we'll rescale it back to the original coordinate system
			pointsVaporized = append(pointsVaporized, point{x: candidatePoint.x + bestPoint.x, y: bestPoint.y - candidatePoint.y})
			delete(distances, candidatePoint)
		}
	}
	return pointsVaporized
}

func visibleNeighbors(origin point, allAsteroids []point) int {
	seen := make(map[string]bool)
	for _, a := range allAsteroids {
		if a == origin {
			continue
		}
		slope := origin.slopeAndDir(a)
		seen[slope.key()] = true
	}
	return len(seen)
}

func main() {
	var asteroids []point
	height := 0
	width := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		width = 0
		line := scanner.Text()
		for index, c := range line {
			if c == '#' {
				asteroids = append(asteroids, point{x: index, y: height})
			}
			width++
		}
		height++
	}

	maxVisibleNeighbors := 0
	var bestPoint point
	for _, a := range asteroids {
		neighbors := visibleNeighbors(a, asteroids)
		if neighbors > maxVisibleNeighbors {
			maxVisibleNeighbors = neighbors
			bestPoint = a
		}
	}
	fmt.Printf("part 1: %d %v\n", maxVisibleNeighbors, bestPoint)

	pointsVaporized := vaporize(bestPoint, asteroids)
	fmt.Printf("part2: %d\n", pointsVaporized[199].x*100+pointsVaporized[199].y)

}
