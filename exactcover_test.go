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

func TestMaxNameLengthEmptyNameSlice(t *testing.T) {
	var ns nameSlice
	if length := ns.maxNameLength(); length != 0 {
		t.Fatalf("expected max name length to be 0 for an empty name slice, but got %d", length)
	}
}

func TestMaxNameLengthNameSliceWithEmptyString(t *testing.T) {
	ns := nameSlice{""}
	if length := ns.maxNameLength(); length != 0 {
		t.Fatalf("expected max name length to be 0 for a name slice containing only an empty string, but got %d", length)
	}
}

func TestMaxNameLengthNameSliceWithDifferentLengthStrings(t *testing.T) {
	ns := nameSlice{"name", "", "longname"}
	if length := ns.maxNameLength(); length != 8 {
		t.Fatalf("expected max name length to be 8 for nameSlice{\"name\", \"\", \"longname\"}, but got %d", length)
	}
}

func TestMaxNameLengthNameSliceWithUnicode(t *testing.T) {
	ns := nameSlice{"name", "", "longname", "⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘"}
	if length := ns.maxNameLength(); length != 12 {
		t.Fatalf("expected max name length to be 8 for nameSlice{\"name\", \"\", \"longname\"}, but got %d", length)
	}
}
