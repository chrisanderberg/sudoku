package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
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

func (ns nameSlice) maxNameLength() int {
	maxLength := 0
	for _, n := range ns {
		if nameLength := utf8.RuneCountInString(string(n)); nameLength > maxLength {
			maxLength = nameLength
		}
	}
	return maxLength
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

func (ec exactCoverDefinition) String() string {
	maxRowNameLength := ec.rowNames.maxNameLength()
	maxColNameLength := ec.colNames.maxNameLength()

	text := make([][]rune, maxColNameLength+len(ec.rowNames)+1)
	for i, _ := range text {
		text[i] = make([]rune, maxRowNameLength+len(ec.colNames)+1)
		for j, _ := range text[i] {
			text[i][j] = ' '
		}
	}

	for i, colName := range ec.colNames {
		offset := maxColNameLength - utf8.RuneCountInString(string(colName))
		for j, character := range colName {
			text[offset+j][maxRowNameLength+1+i] = character
		}
	}

	for i, rowName := range ec.rowNames {
		offset := maxRowNameLength - utf8.RuneCountInString(string(rowName))
		for j, character := range rowName {
			text[maxColNameLength+1+i][offset+j] = character
		}
	}

	for r, _ := range ec.rowNames {
		for c, _ := range ec.colNames {
			if ec.elems[r*len(ec.colNames)+c] {
				text[maxColNameLength+1+r][maxRowNameLength+1+c] = '1'
			} else {
				text[maxColNameLength+1+r][maxRowNameLength+1+c] = '0'
			}
		}
	}

	textLines := make([]string, len(text))
	for i, _ := range text {
		textLines[i] = string(text[i])
	}

	return strings.Join(textLines, "\n")
}

func solve(problem exactCoverDefinition) ([]int, error) {
	if err := problem.validate(); err != nil {
		return nil, err
	}
	return make([]int, 0), nil
}
