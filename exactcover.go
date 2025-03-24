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

// exactCoverConstraints represents the problem definition for an exact cover problem
type exactCoverConstraints struct {
	rowNames nameSlice
	colNames nameSlice
	elems    []bool
}

// validate checks the validity of the exact cover problem definition
func (ec exactCoverConstraints) validate() error {
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
func (ec exactCoverConstraints) String() string {
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

// exactCoverProblem represents a solution to an exact cover problem
type exactCoverProblem struct {
	constraints  exactCoverConstraints
	selectedRows []bool
}

// validate checks the validity of the exact cover solution
func (ecps exactCoverProblem) validate() error {
	if err := ecps.constraints.validate(); err != nil {
		return fmt.Errorf("original problem is invalid: %v", err)
	}
	if len(ecps.selectedRows) != len(ecps.constraints.rowNames) {
		return fmt.Errorf("solution has %d selected rows, but the problem has %d rows", len(ecps.selectedRows), len(ecps.constraints.rowNames))
	}
	coveredCols := make([]string, len(ecps.constraints.colNames))
	for i, isSelected := range ecps.selectedRows {
		if isSelected {
			for j, elem := range ecps.constraints.elems[i*len(ecps.constraints.colNames) : (i+1)*len(ecps.constraints.colNames)] {
				if elem {
					if coveredCols[j] != "" {
						return fmt.Errorf("row %v covers col %v, but col %v is already covered by row %v", ecps.constraints.rowNames[i], ecps.constraints.colNames[j], ecps.constraints.colNames[j], coveredCols[j])
					}
					coveredCols[j] = string(ecps.constraints.rowNames[i])
				}
			}
		}
	}
	return nil
}

func (ecps exactCoverProblem) String() string {
	var selectedRowNames []string
	for i, isSelected := range ecps.selectedRows {
		if isSelected {
			selectedRowNames = append(selectedRowNames, string(ecps.constraints.rowNames[i]))
		}
	}

	return fmt.Sprintf("Selected rows: [%s]\n%s",
		strings.Join(selectedRowNames, ", "),
		ecps.constraints.String())
}

// exactCoverSolution represents a solution to an exact cover problem
type exactCoverSolution exactCoverProblem

// validate checks the validity of the exact cover complete solution
func (eccs exactCoverSolution) validate() error {
	if err := eccs.constraints.validate(); err != nil {
		return fmt.Errorf("original problem is invalid: %v", err)
	}
	if len(eccs.selectedRows) != len(eccs.constraints.rowNames) {
		return fmt.Errorf("solution has %d selected rows, but the problem has %d rows", len(eccs.selectedRows), len(eccs.constraints.rowNames))
	}
	coveredCols := make([]string, len(eccs.constraints.colNames))
	for i, isSelected := range eccs.selectedRows {
		if isSelected {
			for j, elem := range eccs.constraints.elems[i*len(eccs.constraints.colNames) : (i+1)*len(eccs.constraints.colNames)] {
				if elem {
					if coveredCols[j] != "" {
						return fmt.Errorf("row %v covers col %v, but col %v is already covered by row %v", eccs.constraints.rowNames[i], eccs.constraints.colNames[j], eccs.constraints.colNames[j], coveredCols[j])
					}
					coveredCols[j] = string(eccs.constraints.rowNames[i])
				}
			}
		}
	}
	for i, coveredCol := range coveredCols {
		if coveredCol == "" {
			return fmt.Errorf("col %v is not covered by any selected row", eccs.constraints.colNames[i])
		}
	}
	return nil
}

func (eccs exactCoverSolution) String() string {
	return exactCoverProblem(eccs).String()
}

// exactCoverMatrix represents the matrix form of an exact cover problem
type exactCoverMatrix struct {
	constraints exactCoverConstraints
	left        []int
	right       []int
	up          []int
	down        []int
	colSize     []int
	rowNum      []int
	colNum      []int
}

func (ecm *exactCoverMatrix) coverColumn(col int) {
	left := ecm.left
	right := ecm.right
	up := ecm.up
	down := ecm.down
	colNum := ecm.colNum
	colSize := ecm.colSize

	right[left[col]] = right[col]
	left[right[col]] = left[col]
	for i := down[col]; i != col; i = down[i] {
		for j := right[i]; j != i; j = right[j] {
			down[up[j]] = down[j]
			up[down[j]] = up[j]
			colSize[colNum[j]]--
		}
	}
}

func (ecm *exactCoverMatrix) uncoverColumn(col int) {
	left := ecm.left
	right := ecm.right
	up := ecm.up
	down := ecm.down
	colNum := ecm.colNum
	colSize := ecm.colSize

	for i := up[col]; i != col; i = up[i] {
		for j := left[i]; j != i; j = left[j] {
			down[up[j]] = j
			up[down[j]] = j
			colSize[colNum[j]]++
		}
	}
	right[left[col]] = col
	left[right[col]] = col
}

// selectRow covers all columns in the row containing the given element
func (ecm *exactCoverMatrix) selectRow(element int) {
	// Cover each column in the row
	j := element
	for {
		ecm.coverColumn(ecm.colNum[j])
		j = ecm.right[j]
		if j == element {
			break
		}
	}
}

// unselectRow uncovers all columns in the row containing the given element
func (ecm *exactCoverMatrix) unselectRow(element int) {
	// Uncover each column in the row in reverse order
	j := element
	for {
		j = ecm.left[j]
		if j == element {
			break
		}
		ecm.uncoverColumn(ecm.colNum[j])
	}
	ecm.uncoverColumn(ecm.colNum[element])
}

// buildMatrix constructs an exactCoverMatrix from an exactCoverProblem
func buildMatrix(problem exactCoverProblem) (exactCoverMatrix, error) {
	// Validate the problem first
	if err := problem.validate(); err != nil {
		return exactCoverMatrix{}, fmt.Errorf("invalid exact cover problem: %v", err)
	}

	numRows := len(problem.constraints.rowNames)
	numCols := len(problem.constraints.colNames)
	numElems := 0

	// Count the number of true values in problem.constraints.elems
	for _, elem := range problem.constraints.elems {
		if elem {
			numElems++
		}
	}

	// Initialize arrays for the matrix representation
	left := make([]int, numElems+numCols+1)
	right := make([]int, numElems+numCols+1)
	up := make([]int, numElems+numCols+1)
	down := make([]int, numElems+numCols+1)
	rowNum := make([]int, numElems+numCols+1)
	colNum := make([]int, numElems+numCols+1)
	colSize := make([]int, numCols)

	// Set up the column headers
	for i := 0; i < numCols; i++ {
		left[i] = i - 1
		right[i] = i + 1
		up[i] = i
		down[i] = i
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
			if problem.constraints.elems[rowIndex*numCols+colIndex] {
				rowNum[elem] = rowIndex
				colNum[elem] = colIndex

				// insert the element into the column
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

	matrix := exactCoverMatrix{
		constraints: problem.constraints,
		left:        left,
		right:       right,
		up:          up,
		down:        down,
		colSize:     colSize,
		rowNum:      rowNum,
		colNum:      colNum,
	}

	for i, isSelected := range problem.selectedRows {
		if isSelected {
			firstElem, err := matrix.findFirstElementInRow(i)
			if err != nil {
				return exactCoverMatrix{}, err
			}
			matrix.selectRow(firstElem)
		}
	}

	return matrix, nil
}

func (ecm *exactCoverMatrix) findFirstElementInRow(rowNum int) (int, error) {
	for i := len(ecm.constraints.colNames); i < len(ecm.rowNum); i++ {
		if ecm.rowNum[i] == rowNum {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no elements found in row %d", rowNum)
}

// findFirstElementsInRows returns a slice containing the index of the first element in each specified row
func (ecm *exactCoverMatrix) findFirstElementsInRows(rowNums []int) ([]int, error) {
	result := make([]int, 0, len(rowNums))

	// For each row number we're looking for
	for _, targetRow := range rowNums {
		firstElem, err := ecm.findFirstElementInRow(targetRow)
		if err != nil {
			return nil, fmt.Errorf("finding first element in rows: %v", err)
		}
		result = append(result, firstElem)
	}

	return result, nil
}

// findSolution attempts to find a solution to the exact cover problem
// Returns the solution and whether a solution was found
func (ecm *exactCoverMatrix) findSolution() ([]int, bool) {
	// If all columns are covered, we've found a solution
	if ecm.isSolved() {
		return make([]int, 0), true
	}

	// Find column with minimum size
	col := ecm.findSmallestColumn()

	// Try each row in this column
	for i := ecm.down[col]; i != col; i = ecm.down[i] {
		// Add this row to the solution
		ecm.selectRow(i)

		// Recursively search for a solution
		if solution, found := ecm.findSolution(); found {
			return append(solution, ecm.rowNum[i]), true
		}

		// Remove this row from the solution
		ecm.unselectRow(i)
	}

	return nil, false
}

// solve attempts to solve the exact cover problem and returns the solution
func solve(problem exactCoverProblem) (exactCoverSolution, error) {
	if err := problem.validate(); err != nil {
		return exactCoverSolution{}, err
	}

	matrix, err := buildMatrix(problem)
	if err != nil {
		return exactCoverSolution{}, err
	}

	rowNums, found := matrix.findSolution()
	if !found {
		return exactCoverSolution{}, fmt.Errorf("no solution exists")
	}

	// Convert row numbers to selectedRows boolean slice
	selectedRows := make([]bool, len(problem.constraints.rowNames))
	for _, rowNum := range rowNums {
		selectedRows[rowNum] = true
	}

	return exactCoverSolution{
		constraints:  problem.constraints,
		selectedRows: selectedRows,
	}, nil
}

func (ecm *exactCoverMatrix) isSolved() bool {
	var header = len(ecm.left) - 1
	return ecm.right[header] == header
}

// findSmallestColumn returns the column header with the smallest size
func (ecm *exactCoverMatrix) findSmallestColumn() int {
	var header = len(ecm.left) - 1
	minCol := ecm.right[header]
	minSize := ecm.colSize[minCol]

	// Iterate through columns to find smallest
	for colHeader := ecm.right[minCol]; colHeader != header; colHeader = ecm.right[colHeader] {
		if ecm.colSize[colHeader] < minSize {
			minCol = colHeader
			minSize = ecm.colSize[colHeader]
		}
	}

	return minCol
}
