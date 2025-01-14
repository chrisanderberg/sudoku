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
	originalProblem exactCoverProblem
	selectedRows    []bool
}

// validate checks the validity of the exact cover solution
func (ecps exactCoverPartialSolution) validate() error {
	if err := ecps.originalProblem.validate(); err != nil {
		return fmt.Errorf("original problem is invalid: %v", err)
	}
	if len(ecps.selectedRows) != len(ecps.originalProblem.rowNames) {
		return fmt.Errorf("solution has %d selected rows, but the problem has %d rows", len(ecps.selectedRows), len(ecps.originalProblem.rowNames))
	}
	coveredCols := make([]string, len(ecps.originalProblem.colNames))
	for i, isSelected := range ecps.selectedRows {
		if isSelected {
			for j, elem := range ecps.originalProblem.elems[i*len(ecps.originalProblem.colNames) : (i+1)*len(ecps.originalProblem.colNames)] {
				if elem {
					if coveredCols[j] != "" {
						return fmt.Errorf("row %v covers col %v, but col %v is already covered by row %v", ecps.originalProblem.rowNames[i], ecps.originalProblem.colNames[j], ecps.originalProblem.colNames[j], coveredCols[j])
					}
					coveredCols[j] = string(ecps.originalProblem.rowNames[i])
				}
			}
		}
	}
	return nil
}

func (ecps exactCoverPartialSolution) String() string {
	var textLines []string
	for i, rowName := range ecps.originalProblem.rowNames {
		if ecps.selectedRows[i] {
			var colNames []string
			for j, colName := range ecps.originalProblem.colNames {
				if ecps.originalProblem.elems[i*len(ecps.originalProblem.colNames)+j] {
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
	if err := eccs.originalProblem.validate(); err != nil {
		return fmt.Errorf("original problem is invalid: %v", err)
	}
	if len(eccs.selectedRows) != len(eccs.originalProblem.rowNames) {
		return fmt.Errorf("solution has %d selected rows, but the problem has %d rows", len(eccs.selectedRows), len(eccs.originalProblem.rowNames))
	}
	coveredCols := make([]string, len(eccs.originalProblem.colNames))
	for i, isSelected := range eccs.selectedRows {
		if isSelected {
			for j, elem := range eccs.originalProblem.elems[i*len(eccs.originalProblem.colNames) : (i+1)*len(eccs.originalProblem.colNames)] {
				if elem {
					if coveredCols[j] != "" {
						return fmt.Errorf("row %v covers col %v, but col %v is already covered by row %v", eccs.originalProblem.rowNames[i], eccs.originalProblem.colNames[j], eccs.originalProblem.colNames[j], coveredCols[j])
					}
					coveredCols[j] = string(eccs.originalProblem.rowNames[i])
				}
			}
		}
	}
	for i, coveredCol := range coveredCols {
		if coveredCol == "" {
			return fmt.Errorf("col %v is not covered by any selected row", eccs.originalProblem.colNames[i])
		}
	}
	return nil
}

func (eccs exactCoverCompleteSolution) String() string {
	return exactCoverPartialSolution(eccs).String()
}

// exactCoverMatrix represents the matrix form of an exact cover problem
type exactCoverMatrix struct {
	originalProblem exactCoverProblem
}

// solve attempts to solve the exact cover problem and returns a solution
func solve(problem exactCoverProblem) (exactCoverPartialSolution, error) {
	if err := problem.validate(); err != nil {
		return exactCoverPartialSolution{}, err
	}
	return exactCoverPartialSolution{
		originalProblem: problem,
		selectedRows:    make([]bool, len(problem.rowNames)),
	}, nil
}
