package main

import (
	"fmt"
	"log"
)

func main() {
	problem := exactCoverDefinition{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if solution, err := solve(problem); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(solution)
	}
}
