package main

import (
	"errors"
	"fmt"
	"unicode"
)

type name string

func (n name) validate() error {
	for _, r := range n {
		if unicode.IsSpace(r) && r != ' ' {
			return errors.New("name contains invalid whitespace character")
		}
	}
	return nil
}

type nameSlice []name

func (ns nameSlice) validate() error {
	for _, n := range ns {
		if err := n.validate(); err != nil {
			return err
		}
	}
	return nil
}

type exactCoverDefinition struct {
	rowNames nameSlice
	colNames nameSlice
	elems    []bool
}

func (ec exactCoverDefinition) validate() error {
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

/*
func (ec exactCoverDefinition) String() string {

}
*/

func solve(problem exactCoverDefinition) ([]int, error) {
	if err := problem.validate(); err != nil {
		return nil, err
	}
	return make([]int, 0), nil
}
