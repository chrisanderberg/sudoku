package main

import (
	"reflect"
	"testing"
)

func TestValidExactCoverDefinitionHasNoValidationError(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err != nil {
		t.Fatalf("valid exact cover problem shouldn't return an error when validating")
	}
}

func TestNumberRowsValidation(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems with 0 rows should return a validation error")
	}
}

func TestNumberColsValidation(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems with 0 cols should return a validation error")
	}
}

func TestNumberElemsValidation(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems should return a validation error if num elems isn't equal to rows * cols")
	}
}

func TestNamesCantHaveNewlines(t *testing.T) {
	n := name("invalid\nname")
	if err := n.validate(); err == nil {
		t.Fatalf("names with newlines in them should return a validation error")
	}
}

func TestNamesCantHaveCommas(t *testing.T) {
	n := name("invalid,name")
	if err := n.validate(); err == nil {
		t.Fatalf("names with commas in them should return a validation error")
	}
}

func TestPartialSolutionInvalidWhenSameColumnCoveredByMultipleRows(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true, true},
	}
	solution := exactCoverPartialSolution{
		problem:      problem,
		selectedRows: []bool{true, true},
	}
	if err := solution.validate(); err == nil {
		t.Fatalf("partial solution should be invalid when the same column is covered by multiple selected rows")
	}
}

func TestCompleteSolutionInvalidWhenColumnNotCovered(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}
	solution := exactCoverCompleteSolution{
		problem:      problem,
		selectedRows: []bool{true, false},
	}
	if err := solution.validate(); err == nil {
		t.Fatalf("complete solution should be invalid when a column hasn't been covered by any selected row")
	}
}

func TestBuildMatrixIdentity(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 0, 2, 3, 1},
		right:     []int{1, 4, 2, 3, 0},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1, 4},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildMatrixEmpty(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 0, 1},
		right:     []int{1, 2, 0},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1, 2},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildMatrixFull(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 0, 3, 2, 5, 4, 1},
		right:     []int{1, 6, 3, 2, 5, 4, 0},
		up:        []int{4, 5, 0, 1, 2, 3, 6},
		down:      []int{2, 3, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1, 6},
		colSize:   []int{2, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}
