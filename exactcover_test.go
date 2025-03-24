package main

import (
	"reflect"
	"testing"
)

func TestValidExactCoverConstraintsHasNoValidationError(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if err := constraints.validate(); err != nil {
		t.Fatalf("valid exact cover constraints shouldn't return an error when validating")
	}
}

func TestNumberRowsValidation(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if err := constraints.validate(); err == nil {
		t.Fatalf("exact cover constraints with 0 rows should return a validation error")
	}
}

func TestNumberColsValidation(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{},
		elems:    []bool{true},
	}
	if err := constraints.validate(); err == nil {
		t.Fatalf("exact cover constraints with 0 cols should return a validation error")
	}
}

func TestNumberElemsValidation(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if err := constraints.validate(); err == nil {
		t.Fatalf("exact cover constraints should return a validation error if num elems isn't equal to rows * cols")
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

func TestProblemInvalidWhenSameColumnCoveredByMultipleRows(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true, true},
	}
	problem := exactCoverProblem{
		constraints:  constraints,
		selectedRows: []bool{true, true},
	}
	if err := problem.validate(); err == nil {
		t.Fatalf("exact cover problem should be invalid when the same column is covered by multiple selected rows")
	}
}

func TestBuildIdentityMatrix(t *testing.T) {
	problem := exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: nameSlice{"row1", "row2"},
			colNames: nameSlice{"col1", "col2"},
			elems:    []bool{true, false, false, true},
		},
		selectedRows: make([]bool, 2),
	}

	expected := exactCoverMatrix{
		constraints: problem.constraints,
		left:        []int{4, 0, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	actual, err := buildMatrix(problem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildEmptyMatrix(t *testing.T) {
	problem := exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: nameSlice{"row1", "row2"},
			colNames: nameSlice{"col1", "col2"},
			elems:    []bool{false, false, false, false},
		},
		selectedRows: make([]bool, 2),
	}

	expected := exactCoverMatrix{
		constraints: problem.constraints,
		left:        []int{2, 0, 1},
		right:       []int{1, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	actual, err := buildMatrix(problem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildFullMatrix(t *testing.T) {
	problem := exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: nameSlice{"row1", "row2"},
			colNames: nameSlice{"col1", "col2"},
			elems:    []bool{true, true, true, true},
		},
		selectedRows: make([]bool, 2),
	}

	expected := exactCoverMatrix{
		constraints: problem.constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 0},
		up:          []int{4, 5, 0, 1, 2, 3, 6},
		down:        []int{2, 3, 4, 5, 0, 1, 6},
		colSize:     []int{2, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	actual, err := buildMatrix(problem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestBuildInverseIdentityMatrix(t *testing.T) {
	problem := exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: nameSlice{"row1", "row2", "row3"},
			colNames: nameSlice{"col1", "col2", "col3"},
			elems:    []bool{false, true, true, true, false, true, true, true, false},
		},
		selectedRows: make([]bool, 3),
	}

	expected := exactCoverMatrix{
		constraints: problem.constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual, err := buildMatrix(problem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixCoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 0, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 4, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 1},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixUncoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 0, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 4, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 1},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixCoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 0, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 0, 2, 3, 0},
		right:       []int{4, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestIdentityMatrixUncoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, false, false, true},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 0, 2, 3, 1},
		right:       []int{1, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{4, 0, 2, 3, 0},
		right:       []int{4, 4, 2, 3, 0},
		up:          []int{2, 3, 0, 1, 4},
		down:        []int{2, 3, 0, 1, 4},
		colSize:     []int{1, 1},
		rowNum:      []int{-1, -1, 0, 1, -1},
		colNum:      []int{0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixCoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 0, 1},
		right:       []int{1, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 2, 1},
		right:       []int{1, 2, 1},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixUncoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 0, 1},
		right:       []int{1, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 2, 1},
		right:       []int{1, 2, 1},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixCoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 0, 1},
		right:       []int{1, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 0, 0},
		right:       []int{2, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestEmptyMatrixUncoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{false, false, false, false},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 0, 1},
		right:       []int{1, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{2, 0, 0},
		right:       []int{2, 2, 0},
		up:          []int{0, 1, 2},
		down:        []int{0, 1, 2},
		colSize:     []int{0, 0},
		rowNum:      []int{-1, -1, -1},
		colNum:      []int{0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixCoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 0},
		up:          []int{4, 5, 0, 1, 2, 3, 6},
		down:        []int{2, 3, 4, 5, 0, 1, 6},
		colSize:     []int{2, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 6, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 1},
		up:          []int{4, 1, 0, 1, 2, 1, 6},
		down:        []int{2, 1, 4, 5, 0, 1, 6},
		colSize:     []int{2, 0},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixUncoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 0},
		up:          []int{4, 5, 0, 1, 2, 3, 6},
		down:        []int{2, 3, 4, 5, 0, 1, 6},
		colSize:     []int{2, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 6, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 1},
		up:          []int{4, 1, 0, 1, 2, 1, 6},
		down:        []int{2, 1, 4, 5, 0, 1, 6},
		colSize:     []int{2, 0},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixCoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 0},
		up:          []int{4, 5, 0, 1, 2, 3, 6},
		down:        []int{2, 3, 4, 5, 0, 1, 6},
		colSize:     []int{2, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 0},
		right:       []int{6, 6, 3, 2, 5, 4, 0},
		up:          []int{0, 5, 0, 1, 0, 3, 6},
		down:        []int{0, 3, 4, 5, 0, 1, 6},
		colSize:     []int{0, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFullMatrixUncoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2"},
		colNames: nameSlice{"col1", "col2"},
		elems:    []bool{true, true, true, true},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 1},
		right:       []int{1, 6, 3, 2, 5, 4, 0},
		up:          []int{4, 5, 0, 1, 2, 3, 6},
		down:        []int{2, 3, 4, 5, 0, 1, 6},
		colSize:     []int{2, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{6, 0, 3, 2, 5, 4, 0},
		right:       []int{6, 6, 3, 2, 5, 4, 0},
		up:          []int{0, 5, 0, 1, 0, 3, 6},
		down:        []int{0, 3, 4, 5, 0, 1, 6},
		colSize:     []int{0, 2},
		rowNum:      []int{-1, -1, 0, 0, 1, 1, -1},
		colNum:      []int{0, 1, 0, 1, 0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixCoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 9, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 1},
		up:          []int{7, 3, 4, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 1, 2, 7, 2, 0, 1, 9},
		colSize:     []int{2, 1, 1},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.coverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixUncoverLeftColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 9, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 1},
		up:          []int{7, 3, 4, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 1, 2, 7, 2, 0, 1, 9},
		colSize:     []int{2, 1, 1},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.uncoverColumn(0); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixCoverMiddleColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 0, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{2, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{5, 8, 6, 1, 2, 0, 2, 5, 3, 9},
		down:        []int{5, 3, 6, 8, 6, 0, 2, 0, 1, 9},
		colSize:     []int{1, 2, 1},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.coverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixUncoverMiddleColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 0, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{2, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{5, 8, 6, 1, 2, 0, 2, 5, 3, 9},
		down:        []int{5, 3, 6, 8, 6, 0, 2, 0, 1, 9},
		colSize:     []int{1, 2, 1},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.uncoverColumn(1); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixCoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 1},
		right:       []int{1, 9, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 0, 1, 9},
		down:        []int{7, 8, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{1, 1, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.coverColumn(2); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestInverseIdentityMatrixUncoverRightColumn(t *testing.T) {
	constraints := exactCoverConstraints{
		rowNames: nameSlice{"row1", "row2", "row3"},
		colNames: nameSlice{"col1", "col2", "col3"},
		elems:    []bool{false, true, true, true, false, true, true, true, false},
	}

	expected := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 2},
		right:       []int{1, 2, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 5, 3, 9},
		down:        []int{5, 3, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{2, 2, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	actual := exactCoverMatrix{
		constraints: constraints,
		left:        []int{9, 0, 1, 4, 3, 6, 5, 8, 7, 1},
		right:       []int{1, 9, 9, 4, 3, 6, 5, 8, 7, 0},
		up:          []int{7, 8, 6, 1, 2, 0, 4, 0, 1, 9},
		down:        []int{7, 8, 4, 8, 6, 7, 2, 0, 1, 9},
		colSize:     []int{1, 1, 2},
		rowNum:      []int{-1, -1, -1, 0, 0, 1, 1, 2, 2, -1},
		colNum:      []int{0, 1, 2, 1, 2, 0, 2, 0, 1, -1},
	}

	if actual.uncoverColumn(2); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestFindFirstElementsInRow(t *testing.T) {
	problem := exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: nameSlice{"row1", "row2", "row3"},
			colNames: nameSlice{"col1", "col2", "col3"},
			elems:    []bool{false, true, true, true, false, true, true, true, false},
		},
		selectedRows: make([]bool, 3),
	}

	matrix, err := buildMatrix(problem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{3, 5, 7}

	actual, err := matrix.findFirstElementsInRows([]int{0, 1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}
