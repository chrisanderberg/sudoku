package main

import (
	"fmt"
	"log"
)

func main() {
	completeSolution := exactCoverCompleteSolution{
		originalProblem: exactCoverProblem{
			rowNames: nameSlice{"row1", "row2", "row3"},
			colNames: nameSlice{"col1", "col2", "col3"},
			elems:    []bool{true, true, false, false, true, false, false, false, true},
		},
		selectedRows: []bool{true, true, true},
	}
	if err := completeSolution.validate(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(completeSolution)
	}
}
