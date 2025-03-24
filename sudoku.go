package main

import (
	"fmt"
	"os"
)

func main() {
	problem := exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: nameSlice{"row1", "row2", "row3"},
			colNames: nameSlice{"col1", "col2", "col3"},
			elems:    []bool{true, true, false, false, true, false, false, false, true},
		},
		selectedRows: []bool{false, true, false},
	}
	fmt.Println(problem)
	solution, err := solve(problem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving problem: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(solution)
}
