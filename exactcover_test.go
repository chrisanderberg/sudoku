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

func TestBuildIdentityMatrix(t *testing.T) {
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
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildEmptyMatrix(t *testing.T) {
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
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildFullMatrix(t *testing.T) {
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
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildInverseIdentityMatrix(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual := buildMatrix(problem); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixCoverLeftColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 0, 2, 3, 1},
		right:     []int{1, 4, 2, 3, 0},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 4, 2, 3, 1},
		right:     []int{1, 4, 2, 3, 1},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixUncoverLeftColumn(t *testing.T) {
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
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 4, 2, 3, 1},
		right:     []int{1, 4, 2, 3, 1},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixCoverRightColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 0, 2, 3, 1},
		right:     []int{1, 4, 2, 3, 0},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 0, 2, 3, 0},
		right:     []int{4, 4, 2, 3, 0},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixUncoverRightColumn(t *testing.T) {
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
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{4, 0, 2, 3, 0},
		right:     []int{4, 4, 2, 3, 0},
		up:        []int{2, 3, 0, 1, 4},
		down:      []int{2, 3, 0, 1, 4},
		colHeader: []int{0, 1, 0, 1},
		colSize:   []int{1, 1},
		rowNum:    []int{-1, -1, 0, 1, -1},
		colNum:    []int{0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixCoverLeftColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 0, 1},
		right:     []int{1, 2, 0},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 2, 1},
		right:     []int{1, 2, 1},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixUnoverLeftColumn(t *testing.T) {
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
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 2, 1},
		right:     []int{1, 2, 1},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixCoverRightColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 0, 1},
		right:     []int{1, 2, 0},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 0, 0},
		right:     []int{2, 2, 0},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixUncoverRightColumn(t *testing.T) {
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
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{2, 0, 0},
		right:     []int{2, 2, 0},
		up:        []int{0, 1, 2},
		down:      []int{0, 1, 2},
		colHeader: []int{0, 1},
		colSize:   []int{0, 0},
		rowNum:    []int{-1, -1, -1},
		colNum:    []int{0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixCoverLeftColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 0, 3, 2, 5, 4, 1},
		right:     []int{1, 6, 3, 2, 5, 4, 0},
		up:        []int{4, 5, 0, 1, 2, 3, 6},
		down:      []int{2, 3, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 6, 3, 2, 5, 4, 1},
		right:     []int{1, 6, 3, 2, 5, 4, 1},
		up:        []int{4, 1, 0, 1, 2, 1, 6},
		down:      []int{2, 1, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 0},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixUncoverLeftColumn(t *testing.T) {
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
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 6, 3, 2, 5, 4, 1},
		right:     []int{1, 6, 3, 2, 5, 4, 1},
		up:        []int{4, 1, 0, 1, 2, 1, 6},
		down:      []int{2, 1, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 0},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixCoverRightColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 0, 3, 2, 5, 4, 1},
		right:     []int{1, 6, 3, 2, 5, 4, 0},
		up:        []int{4, 5, 0, 1, 2, 3, 6},
		down:      []int{2, 3, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 0, 3, 2, 5, 4, 0},
		right:     []int{6, 6, 3, 2, 5, 4, 0},
		up:        []int{0, 5, 0, 1, 0, 3, 6},
		down:      []int{0, 3, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{0, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixUnoverRightColumn(t *testing.T) {
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
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{2, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{6, 0, 3, 2, 5, 4, 0},
		right:     []int{6, 6, 3, 2, 5, 4, 0},
		up:        []int{0, 5, 0, 1, 0, 3, 6},
		down:      []int{0, 3, 4, 5, 0, 1, 6},
		colHeader: []int{0, 1, 0, 1, 0, 1},
		colSize:   []int{0, 2},
		rowNum:    []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:    []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixCoverLeftColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 9, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 1},
		up:        []int{7, 3, 4, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 1, 2, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 1, 1},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixUncoverLeftColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 9, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 1},
		up:        []int{7, 3, 4, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 1, 2, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 1, 1},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixCoverMiddleColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 0, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{2, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{5, 8, 6, 1, 2, 0, 2, 5, 3, 9},
		down:      []int{5, 3, 6, 8, 6, 0, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{1, 2, 1},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixUncoverMiddleColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 0, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{2, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{5, 8, 6, 1, 2, 0, 2, 5, 3, 9},
		down:      []int{5, 3, 6, 8, 6, 0, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{1, 2, 1},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixCoverRightColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 1},
		right:     []int{1, 9, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 0, 1, 9},
		down:      []int{7, 8, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{1, 1, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.coverColumn(2); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixUncoverRightColumn(t *testing.T) {
	problem := exactCoverProblem{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:     []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:      []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{2, 2, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		problem:   problem,
		left:      []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 1},
		right:     []int{1, 9, 9, 4, 3, 6, 5, 8, 7, 0},
		up:        []int{7, 8, 6, 1, 2, 0, 4, 0, 1, 9},
		down:      []int{7, 8, 4, 8, 6, 7, 2, 0, 1, 9},
		colHeader: []int{0, 1, 2, 1, 2, 0, 2, 0, 1},
		colSize:   []int{1, 1, 2},
		rowNum:    []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:    []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.uncoverColumn(2); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}
