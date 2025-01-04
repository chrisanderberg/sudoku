package main

import (
	"testing"
)

func TestNumberRowsValidation(t *testing.T) {
	problem := exactCoverDefinition{
		rowNames: []string{},
		colNames: []string{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems with 0 rows should return a validation error")
	}
}

func TestNumberColsValidation(t *testing.T) {
	problem := exactCoverDefinition{
		rowNames: []string{"row1"},
		colNames: []string{},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems with 0 cols should return a validation error")
	}
}

func TestNumberElemsValidation(t *testing.T) {
	problem := exactCoverDefinition{
		rowNames: []string{"row1", "row2"},
		colNames: []string{"col1"},
		elems:    []bool{true},
	}
	if _, err := solve(problem); err == nil {
		t.Fatalf("exact cover problems should return a validation error if num elems isn't equal to rows * cols")
	}
}
