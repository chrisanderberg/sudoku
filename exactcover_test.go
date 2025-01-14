package main

import (
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
		originalProblem: problem,
		selectedRows:    []bool{true, true},
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
		originalProblem: problem,
		selectedRows:    []bool{true, false},
	}
	if err := solution.validate(); err == nil {
		t.Fatalf("complete solution should be invalid when a column hasn't been covered by any selected row")
	}
}
