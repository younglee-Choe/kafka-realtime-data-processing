package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	var indexArr [][]int
	// combin provides several ways to work with the combinations of
	// different objects. CombinationGenerator constructs an iterator
	// for the combinations.
	n := 4
	k := n - 1
	gen := combin.NewCombinationGenerator(n, k)
	// idx := 0
	for gen.Next() {
		fmt.Println(gen.Combination(nil)) // can also store in-place.
		indexArr = append(indexArr, gen.Combination(nil))
		// idx++
	}
	fmt.Println("indexArr:", indexArr[0], indexArr[1], indexArr[2], indexArr[3])

	// Output:
	// 0 [0 1 2]
	// 1 [0 1 3]
	// 2 [0 2 3]
	// 3 [1 2 3]
}
