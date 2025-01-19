package main

import (
	"fmt"
)

func main() {
	completeSolution := exactCoverCompleteSolution{
		problem: exactCoverProblem{
			rowNames: nameSlice{"row1", "row2", "row3"},
			colNames: nameSlice{"col1", "col2", "col3"},
			elems:    []bool{true, true, false, false, true, false, false, false, true},
		},
		selectedRows: []bool{true, true, true},
	}

	fmt.Println(buildMatrix(completeSolution.problem))
}
