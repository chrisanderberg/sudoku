package main

import (
	"fmt"
	"log"
)

func main() {
	problem := exactCoverDefinition{
		rowNames: []string{"row1"},
		colNames: []string{"col1"},
		elems:    []bool{true},
	}
	if solution, err := solve(problem); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(solution)
	}
}
