package main

import (
	"errors"
	"fmt"
)

type exactCoverDefinition struct {
	rowNames []string
	colNames []string
	elems    []bool
}

func (ec exactCoverDefinition) validate() error {
	if len(ec.rowNames) < 1 {
		return errors.New(fmt.Sprintf("exact cover must have at least 1 row, but %d rows were provided", len(ec.rowNames)))
	}
	if len(ec.colNames) < 1 {
		return errors.New(fmt.Sprintf("exact cover must have at least 1 col, but %d cols were provided", len(ec.colNames)))
	}
	if len(ec.elems) != len(ec.rowNames)*len(ec.colNames) {
		return errors.New(fmt.Sprintf("exact cover with %d rows and %d cols should have %d*%d=%d elems, but %d elems were provided instead", len(ec.rowNames), len(ec.colNames), len(ec.rowNames), len(ec.colNames), len(ec.rowNames)*len(ec.colNames), len(ec.elems)))
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
