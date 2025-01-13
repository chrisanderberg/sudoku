package main

import (
	"fmt"
	"log"
)

func main() {
	problem := exactCoverDefinition{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{true, false, false, false, true, false, false, false, true},
	}
	if solution, err := solve(problem); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(solution)
	}
	fmt.Println(problem)
}
