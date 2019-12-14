package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

type component struct {
	count    int64
	chemical string
}

func parse(s string) component {
	var c = component{}
	fmt.Sscanf(s, "%d %s", &c.count, &c.chemical)
	return c
}

type reaction struct {
	inputs []component
	output component
}

var oreCollected int64

func consume(c component, inventory map[string]int64, reactions map[string]reaction) {
	if c.chemical == "ORE" {
		oreCollected += c.count
	} else {
		need := int64(0)
		if inventory[c.chemical] < c.count {
			need = c.count - inventory[c.chemical]
			produce(component{chemical: c.chemical, count: need}, inventory, reactions)
		}
		inventory[c.chemical] = inventory[c.chemical] - c.count
	}
}

func produce(c component, inventory map[string]int64, reactions map[string]reaction) {
	r := reactions[c.chemical]
	iterations := int64(math.Ceil(float64(c.count) / float64(r.output.count)))
	for _, c := range r.inputs {
		consume(component{chemical: c.chemical, count: c.count * iterations}, inventory, reactions)
	}
	inventory[c.chemical] = inventory[c.chemical] + r.output.count*iterations
}

func main() {

	reactions := make(map[string]reaction)
	scanner := bufio.NewScanner(os.Stdin)
	var line = regexp.MustCompile(`^(.*) => (.*)$`)

	for scanner.Scan() {
		r := reaction{}
		matches := line.FindStringSubmatch(scanner.Text())
		r.output = parse(matches[2])
		for _, s := range strings.Split(matches[1], ", ") {
			r.inputs = append(r.inputs, parse(s))
		}
		reactions[r.output.chemical] = r
	}

	/// part 1
	{
		var inventory = make(map[string]int64)
		produce(component{chemical: "FUEL", count: 1}, inventory, reactions)
		fmt.Printf("part 1: %d\n", oreCollected)
	}

	/// part 2
	{
		best := sort.Search(100000000, func(i int) bool {
			oreCollected = 0
			var inventory = make(map[string]int64)
			produce(component{chemical: "FUEL", count: int64(i)}, inventory, reactions)
			return oreCollected > 1000000000000
		})
		/// -1 as search returns [0,n]
		fmt.Printf("part 2:%d\n", best-1)
	}
}
