package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func countChildren(nodes map[string][]string, child string) int {
	count := 0
	if children, ok := nodes[child]; ok != false {
		for _, c := range children {
			count = count + countChildren(nodes, c) + 1
		}
	}
	return count
}

func findParents(child2Parent map[string]string, child string) []string {
	var parents []string
	for {
		if parent, ok := child2Parent[child]; ok == false {
			break
		} else {
			parents = append(parents, parent)
			child = parent
		}
	}
	return parents
}

func main() {

	child2Parent := make(map[string]string)
	nodes := make(map[string][]string)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pieces := strings.Split(scanner.Text(), ")")
		parent, child := pieces[0], pieces[1]

		if _, ok := nodes[parent]; ok != true {
			nodes[parent] = []string{child}
		} else {
			nodes[parent] = append(nodes[parent], child)
		}
		child2Parent[child] = parent
	}

	// part 1
	orbitCount := 0
	for _, children := range nodes {
		for _, c := range children {
			// + 1 to include the link between this node and the child itself
			orbitCount += countChildren(nodes, c) + 1
		}
	}
	fmt.Printf("orbit count: %d\n", orbitCount)

	// part 2
	myOrbitsParents := findParents(child2Parent, child2Parent["YOU"])
	santasOrbitsParents := findParents(child2Parent, child2Parent["SAN"])

	// for quick lookup when finding closest ancestor
	santaMap := make(map[string]bool)
	for _, v := range santasOrbitsParents {
		santaMap[v] = true
	}

	myOrbitToClosestSharedParent := 0
	var closestSharedParent string
	for _, value := range myOrbitsParents {
		myOrbitToClosestSharedParent++
		if _, ok := santaMap[value]; ok != false {
			closestSharedParent = value
			break
		}
	}

	santasOrbitToClosestSharedParent := 0
	for _, value := range santasOrbitsParents {
		santasOrbitToClosestSharedParent++
		if value == closestSharedParent {
			break
		}
	}

	fmt.Printf("transfer count: %d\n", santasOrbitToClosestSharedParent+myOrbitToClosestSharedParent)
}
