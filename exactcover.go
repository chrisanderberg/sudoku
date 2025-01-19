package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// name represents a string type with validation for whitespace characters
type name string

// validate checks that the only whitespace characters in the name are spaces and that it doesn't contain commas
func (n name) validate() error {
	for _, r := range n {
		if unicode.IsSpace(r) && r != ' ' {
			return errors.New("name contains invalid whitespace character")
		}
		if r == ',' {
			return errors.New("name contains invalid comma character")
		}
	}
	return nil
}

// nameSlice represents a slice of names with validation and utility methods
type nameSlice []name

// validate checks each name in the slice for validity
func (ns nameSlice) validate() error {
	for _, n := range ns {
		if err := n.validate(); err != nil {
			return err
		}
	}
	return nil
}

// exactCoverProblem represents the problem definition for an exact cover problem
type exactCoverProblem struct {
	rowNames nameSlice
	colNames nameSlice
	elems    []bool
}

// validate checks the validity of the exact cover problem definition
func (ec exactCoverProblem) validate() error {
	if len(ec.rowNames) < 1 {
		return fmt.Errorf("exact cover must have at least 1 row, but %d rows were provided", len(ec.rowNames))
	}
	if len(ec.colNames) < 1 {
		return fmt.Errorf("exact cover must have at least 1 col, but %d cols were provided", len(ec.colNames))
	}
	if len(ec.elems) != len(ec.rowNames)*len(ec.colNames) {
		return fmt.Errorf("exact cover with %d rows and %d cols should have %d*%d=%d elems, but %d elems were provided instead", len(ec.rowNames), len(ec.colNames), len(ec.rowNames), len(ec.colNames), len(ec.rowNames)*len(ec.colNames), len(ec.elems))
	}
	if err := ec.rowNames.validate(); err != nil {
		return fmt.Errorf("invalid row names: %v", err)
	}
	if err := ec.colNames.validate(); err != nil {
		return fmt.Errorf("invalid col names: %v", err)
	}
	return nil
}

// String returns a string representation of the exact cover problem
func (ec exactCoverProblem) String() string {
	var textLines []string
	for i, rowName := range ec.rowNames {
		var colNames []string
		for j, colName := range ec.colNames {
			if ec.elems[i*len(ec.colNames)+j] {
				colNames = append(colNames, string(colName))
			}
		}
		textLines = append(textLines, string(rowName)+": "+strings.Join(colNames, ", "))
	}

	return strings.Join(textLines, "\n")
}

// exactCoverPartialSolution represents a solution to an exact cover problem
type exactCoverPartialSolution struct {
	problem      exactCoverProblem
	selectedRows []bool
}

// validate checks the validity of the exact cover solution
func (ecps exactCoverPartialSolution) validate() error {
	if err := ecps.problem.validate(); err != nil {
		return fmt.Errorf("original problem is invalid: %v", err)
	}
	if len(ecps.selectedRows) != len(ecps.problem.rowNames) {
		return fmt.Errorf("solution has %d selected rows, but the problem has %d rows", len(ecps.selectedRows), len(ecps.problem.rowNames))
	}
	coveredCols := make([]string, len(ecps.problem.colNames))
	for i, isSelected := range ecps.selectedRows {
		if isSelected {
			for j, elem := range ecps.problem.elems[i*len(ecps.problem.colNames) : (i+1)*len(ecps.problem.colNames)] {
				if elem {
					if coveredCols[j] != "" {
						return fmt.Errorf("row %v covers col %v, but col %v is already covered by row %v", ecps.problem.rowNames[i], ecps.problem.colNames[j], ecps.problem.colNames[j], coveredCols[j])
					}
					coveredCols[j] = string(ecps.problem.rowNames[i])
				}
			}
		}
	}
	return nil
}

func (ecps exactCoverPartialSolution) String() string {
	var textLines []string
	for i, rowName := range ecps.problem.rowNames {
		if ecps.selectedRows[i] {
			var colNames []string
			for j, colName := range ecps.problem.colNames {
				if ecps.problem.elems[i*len(ecps.problem.colNames)+j] {
					colNames = append(colNames, string(colName))
				}
			}
			textLines = append(textLines, string(rowName)+": "+strings.Join(colNames, ", "))
		}
	}

	return strings.Join(textLines, "\n")
}

type exactCoverCompleteSolution exactCoverPartialSolution

// validate checks the validity of the exact cover complete solution
func (eccs exactCoverCompleteSolution) validate() error {
	if err := eccs.problem.validate(); err != nil {
		return fmt.Errorf("original problem is invalid: %v", err)
	}
	if len(eccs.selectedRows) != len(eccs.problem.rowNames) {
		return fmt.Errorf("solution has %d selected rows, but the problem has %d rows", len(eccs.selectedRows), len(eccs.problem.rowNames))
	}
	coveredCols := make([]string, len(eccs.problem.colNames))
	for i, isSelected := range eccs.selectedRows {
		if isSelected {
			for j, elem := range eccs.problem.elems[i*len(eccs.problem.colNames) : (i+1)*len(eccs.problem.colNames)] {
				if elem {
					if coveredCols[j] != "" {
						return fmt.Errorf("row %v covers col %v, but col %v is already covered by row %v", eccs.problem.rowNames[i], eccs.problem.colNames[j], eccs.problem.colNames[j], coveredCols[j])
					}
					coveredCols[j] = string(eccs.problem.rowNames[i])
				}
			}
		}
	}
	for i, coveredCol := range coveredCols {
		if coveredCol == "" {
			return fmt.Errorf("col %v is not covered by any selected row", eccs.problem.colNames[i])
		}
	}
	return nil
}

func (eccs exactCoverCompleteSolution) String() string {
	return exactCoverPartialSolution(eccs).String()
}

// exactCoverMatrix represents the matrix form of an exact cover problem
type exactCoverMatrix struct {
	problem   exactCoverProblem
	left      []int
	right     []int
	up        []int
	down      []int
	colHeader []int
	colSize   []int
	rowNum    []int
	colNum    []int
}

func (ecm exactCoverMatrix) coverColumn(col int) {
	left := ecm.left
	right := ecm.right
	up := ecm.up
	down := ecm.down
	colHeader := ecm.colHeader
	colSize := ecm.colSize

	right[left[col]] = right[col]
	left[right[col]] = left[col]
	for i := down[col]; i != col; i = down[i] {
		for j := right[i]; j != i; j = right[j] {
			down[up[j]] = down[j]
			up[down[j]] = up[j]
			colSize[colHeader[j]]--
		}
	}
}

func (ecm exactCoverMatrix) uncoverColumn(col int) {
	left := ecm.left
	right := ecm.right
	up := ecm.up
	down := ecm.down
	colHeader := ecm.colHeader
	colSize := ecm.colSize

	for i := up[col]; i != col; i = up[i] {
		for j := left[i]; j != i; j = left[j] {
			down[up[j]] = j
			up[down[j]] = j
			colSize[colHeader[j]]++
		}
	}
	right[left[col]] = col
	left[right[col]] = col
}

// buildMatrix constructs an exactCoverMatrix from an exactCoverPartialSolution
func buildMatrix(problem exactCoverProblem) exactCoverMatrix {
	numRows := len(problem.rowNames)
	numCols := len(problem.colNames)
	numElems := 0

	// Count the number of true values in problem.elems
	for _, elem := range problem.elems {
		if elem {
			numElems++
		}
	}

	// Initialize arrays for the matrix representation
	left := make([]int, numElems+numCols+1)
	right := make([]int, numElems+numCols+1)
	up := make([]int, numElems+numCols+1)
	down := make([]int, numElems+numCols+1)
	colHeader := make([]int, numElems+numCols)
	rowNum := make([]int, numElems+numCols+1)
	colNum := make([]int, numElems+numCols+1)
	colSize := make([]int, numCols)

	// Set up the column headers
	for i := 0; i < numCols; i++ {
		left[i] = i - 1
		right[i] = i + 1
		up[i] = i
		down[i] = i
		colHeader[i] = i
		rowNum[i] = -1
		colNum[i] = i
		colSize[i] = 0
	}

	// Set up the header node
	header := numElems + numCols
	left[0] = header
	left[header] = numCols - 1
	right[header] = 0
	right[numCols-1] = header
	up[header] = header
	down[header] = header
	colNum[header] = -1
	rowNum[header] = -1

	// Populate the matrix with elements
	elem := numCols
	for rowIndex := 0; rowIndex < numRows; rowIndex++ {
		firstElemInRow := true
		for colIndex := 0; colIndex < numCols; colIndex++ {
			if problem.elems[rowIndex*numCols+colIndex] {
				rowNum[elem] = rowIndex
				colNum[elem] = colIndex

				// insert the element into the column
				colHeader[elem] = colIndex
				colSize[colIndex]++
				down[elem] = colIndex
				up[elem] = up[colIndex]
				down[up[colIndex]] = elem
				up[colIndex] = elem

				if firstElemInRow {
					// first element in the row is its own left and right
					firstElemInRow = false
					left[elem] = elem
					right[elem] = elem
				} else {
					// insert the element into the row
					left[elem] = elem - 1
					right[elem] = right[elem-1]
					left[right[elem]] = elem
					right[elem-1] = elem
				}

				elem++
			}
		}
	}

	return exactCoverMatrix{
		problem:   problem,
		left:      left,
		right:     right,
		up:        up,
		down:      down,
		colHeader: colHeader,
		colSize:   colSize,
		rowNum:    rowNum,
		colNum:    colNum,
	}
}

// solve attempts to solve the exact cover problem and returns a solution
func solve(problem exactCoverProblem) (exactCoverPartialSolution, error) {
	if err := problem.validate(); err != nil {
		return exactCoverPartialSolution{}, err
	}
	return exactCoverPartialSolution{
		problem:      problem,
		selectedRows: make([]bool, len(problem.rowNames)),
	}, nil
}
