package main

import (
	"testing"
)

func TestValidExactCoverDefinitionHasNoValidationError(t *testing.T) {
	problem := exactCoverDefinition{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err != nil {
		t.Fatalf("valid exact cover problem shouldn't return an error when validating")
	}
}

func TestNumberRowsValidation(t *testing.T) {
	problem := exactCoverDefinition{
		rowNames: nameSlice{},
		colNames: nameSlice{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems with 0 rows should return a validation error")
	}
}

func TestNumberColsValidation(t *testing.T) {
	problem := exactCoverDefinition{
		rowNames: nameSlice{"row1"},
		colNames: nameSlice{},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems with 0 cols should return a validation error")
	}
}

func TestNumberElemsValidation(t *testing.T) {
	problem := exactCoverDefinition{
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
