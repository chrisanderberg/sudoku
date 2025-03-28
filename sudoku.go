package main

import (
	"fmt"
	"os"
	"strings"
)

// sudokuProblem represents a 9x9 Sudoku grid where 0 represents an empty cell
type sudokuProblem struct {
	cells [81]int // 9x9 grid stored as a flat array
}

// validate checks that the Sudoku problem is valid
func (sp sudokuProblem) validate() error {
	// Check each cell is in valid range (0-9)
	for i, val := range sp.cells {
		if val < 0 || val > 9 {
			return fmt.Errorf("cell %d contains invalid value %d: must be between 0 and 9", i, val)
		}
	}

	// Check no duplicate numbers in each row
	for row := 0; row < 9; row++ {
		seen := make(map[int]bool)
		for col := 0; col < 9; col++ {
			val := sp.cells[row*9+col]
			if val != 0 {
				if seen[val] {
					return fmt.Errorf("row %d contains duplicate value %d", row+1, val)
				}
				seen[val] = true
			}
		}
	}

	// Check no duplicate numbers in each column
	for col := 0; col < 9; col++ {
		seen := make(map[int]bool)
		for row := 0; row < 9; row++ {
			val := sp.cells[row*9+col]
			if val != 0 {
				if seen[val] {
					return fmt.Errorf("column %d contains duplicate value %d", col+1, val)
				}
				seen[val] = true
			}
		}
	}

	// Check no duplicate numbers in each 3x3 box
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			seen := make(map[int]bool)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					row := boxRow*3 + i
					col := boxCol*3 + j
					val := sp.cells[row*9+col]
					if val != 0 {
						if seen[val] {
							return fmt.Errorf("3x3 box at position (%d,%d) contains duplicate value %d", boxRow+1, boxCol+1, val)
						}
						seen[val] = true
					}
				}
			}
		}
	}

	return nil
}

// String returns a string representation of the Sudoku grid
func (sp sudokuProblem) String() string {
	var sb strings.Builder
	sb.WriteString("┌───────┬───────┬───────┐\n")

	for row := 0; row < 9; row++ {
		if row > 0 && row%3 == 0 {
			sb.WriteString("├───────┼───────┼───────┤\n")
		}
		sb.WriteString("│")
		for col := 0; col < 9; col++ {
			if col > 0 && col%3 == 0 {
				sb.WriteString(" │")
			}
			val := sp.cells[row*9+col]
			if val == 0 {
				sb.WriteString(" ·")
			} else {
				sb.WriteString(fmt.Sprintf(" %d", val))
			}
		}
		sb.WriteString(" │\n")
	}

	sb.WriteString("└───────┴───────┴───────┘\n")
	return sb.String()
}

// sudokuSolution represents a complete solution to a Sudoku puzzle
type sudokuSolution struct {
	cells [81]int // 9x9 grid stored as a flat array
}

// validate checks that the Sudoku solution is valid and complete
func (ss sudokuSolution) validate() error {
	// First check all basic Sudoku rules using the problem validator
	if err := sudokuProblem(ss).validate(); err != nil {
		return err
	}

	// Then verify that there are no empty cells
	for i, val := range ss.cells {
		if val == 0 {
			return fmt.Errorf("cell %d is empty", i)
		}
	}

	return nil
}

// String returns a string representation of the solved Sudoku grid
func (ss sudokuSolution) String() string {
	return sudokuProblem(ss).String()
}

// toExactCover converts a Sudoku problem to an exact cover problem
func (sp sudokuProblem) toExactCover() exactCoverProblem {
	// For a 9x9 Sudoku:
	// - 9x9=81 cell constraints (one number per cell)
	// - 9x9=81 row constraints (each number appears once per row)
	// - 9x9=81 column constraints (each number appears once per column)
	// - 9x9=81 box constraints (each number appears once per 3x3 box)
	// Total: 324 constraints

	// Each row in the exact cover matrix represents placing a specific number
	// in a specific cell. For a 9x9 grid, there are 9x9x9=729 possibilities.

	// Create row names (e.g., "R1C2#3" means "put 3 in row 1, column 2")
	rowNames := make(nameSlice, 729)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			for num := 1; num <= 9; num++ {
				idx := row*81 + col*9 + (num - 1)
				rowNames[idx] = name(fmt.Sprintf("R%dC%d#%d", row+1, col+1, num))
			}
		}
	}

	// Create column names for each constraint type
	colNames := make(nameSlice, 324)
	idx := 0
	// Cell constraints (e.g., "R1C2" means "cell at row 1, column 2 must have one number")
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			colNames[idx] = name(fmt.Sprintf("R%dC%d", row+1, col+1))
			idx++
		}
	}
	// Row constraints (e.g., "R1#3" means "row 1 must have one 3")
	for row := 0; row < 9; row++ {
		for num := 1; num <= 9; num++ {
			colNames[idx] = name(fmt.Sprintf("R%d#%d", row+1, num))
			idx++
		}
	}
	// Column constraints (e.g., "C1#3" means "column 1 must have one 3")
	for col := 0; col < 9; col++ {
		for num := 1; num <= 9; num++ {
			colNames[idx] = name(fmt.Sprintf("C%d#%d", col+1, num))
			idx++
		}
	}
	// Box constraints (e.g., "B1#3" means "box 1 must have one 3")
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			for num := 1; num <= 9; num++ {
				colNames[idx] = name(fmt.Sprintf("B%d#%d", boxRow*3+boxCol+1, num))
				idx++
			}
		}
	}

	// Create the constraint matrix
	elems := make([]bool, 729*324)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			for num := 1; num <= 9; num++ {
				rowIdx := row*81 + col*9 + (num - 1)

				// Cell constraint
				cellIdx := row*9 + col
				elems[rowIdx*324+cellIdx] = true

				// Row constraint
				rowConstraintIdx := 81 + row*9 + (num - 1)
				elems[rowIdx*324+rowConstraintIdx] = true

				// Column constraint
				colConstraintIdx := 162 + col*9 + (num - 1)
				elems[rowIdx*324+colConstraintIdx] = true

				// Box constraint
				boxRow, boxCol := row/3, col/3
				boxIdx := boxRow*3 + boxCol
				boxConstraintIdx := 243 + boxIdx*9 + (num - 1)
				elems[rowIdx*324+boxConstraintIdx] = true
			}
		}
	}

	// Create selectedRows slice based on the given numbers in the Sudoku grid
	selectedRows := make([]bool, 729)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			num := sp.cells[row*9+col]
			if num != 0 {
				rowIdx := row*81 + col*9 + (num - 1)
				selectedRows[rowIdx] = true
			}
		}
	}

	return exactCoverProblem{
		constraints: exactCoverConstraints{
			rowNames: rowNames,
			colNames: colNames,
			elems:    elems,
		},
		selectedRows: selectedRows,
	}
}

// fromExactCover converts an exact cover solution back to a Sudoku solution
func fromExactCover(ecs exactCoverSolution) sudokuSolution {
	var solution sudokuSolution

	// Each selected row in the exact cover solution represents a number placement
	// The row name format is "R{row}C{col}#{num}"
	for i, isSelected := range ecs.selectedRows {
		if isSelected {
			row := (i / 81)
			col := (i / 9) % 9
			num := (i % 9) + 1
			solution.cells[row*9+col] = num
		}
	}

	return solution
}

// solve converts the Sudoku problem to an exact cover problem, solves it,
// and converts the solution back to a Sudoku solution
func (sp sudokuProblem) solve() (sudokuSolution, error) {
	if err := sp.validate(); err != nil {
		return sudokuSolution{}, fmt.Errorf("invalid sudoku problem: %v", err)
	}

	exactCoverProblem := sp.toExactCover()
	exactCoverSolution, err := solve(exactCoverProblem)
	if err != nil {
		return sudokuSolution{}, fmt.Errorf("solving sudoku: %v", err)
	}

	solution := fromExactCover(exactCoverSolution)
	if err := solution.validate(); err != nil {
		return sudokuSolution{}, fmt.Errorf("invalid solution: %v", err)
	}

	return solution, nil
}

// fromString creates a sudokuProblem from a string representation.
// The string should contain 81 characters where:
// - digits 1-9 represent filled cells
// - '.' or '0' represent empty cells
// - whitespace and other characters are ignored
func fromString(s string) (sudokuProblem, error) {
	var problem sudokuProblem
	cellIndex := 0

	// Process each character
	for _, r := range s {
		// Skip whitespace and other formatting characters
		if strings.ContainsRune(" \t\n\r│─┌┐└┘├┤┬┴┼", r) {
			continue
		}

		// Check if we've already filled all cells
		if cellIndex >= 81 {
			return sudokuProblem{}, fmt.Errorf("input contains more than 81 cells")
		}

		// Parse the cell value
		switch r {
		case '.', '·', '0':
			problem.cells[cellIndex] = 0
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			problem.cells[cellIndex] = int(r - '0')
		default:
			return sudokuProblem{}, fmt.Errorf("invalid character '%c' at position %d: must be 0-9 or '.'", r, cellIndex)
		}
		cellIndex++
	}

	// Check if we got enough cells
	if cellIndex < 81 {
		return sudokuProblem{}, fmt.Errorf("input contains only %d cells, expected 81", cellIndex)
	}

	// Validate the problem
	if err := problem.validate(); err != nil {
		return sudokuProblem{}, fmt.Errorf("invalid sudoku problem: %v", err)
	}

	return problem, nil
}

func main() {
	// Check if we have enough arguments
	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> <puzzle-file> [output-file]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  display  Display the formatted puzzle\n")
		fmt.Fprintf(os.Stderr, "  format   Format the puzzle file in-place\n")
		fmt.Fprintf(os.Stderr, "  solve    Solve the puzzle and display the solution\n")
		fmt.Fprintf(os.Stderr, "          (optionally write solution to output-file)\n")
		os.Exit(1)
	}

	command := os.Args[1]
	filepath := os.Args[2]

	// Validate command and argument count
	switch command {
	case "solve":
		// solve can have 1 or 2 file arguments
		if len(os.Args) > 4 {
			fmt.Fprintf(os.Stderr, "solve command takes at most 2 file arguments\n")
			os.Exit(1)
		}
	default:
		// other commands must have exactly 1 file argument
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "%s command takes exactly 1 file argument\n", command)
			os.Exit(1)
		}
	}

	// Read the file contents
	input, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the puzzle
	problem, err := fromString(string(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing problem: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "display":
		fmt.Println(problem)
	case "format":
		// Write the formatted puzzle back to the file
		formatted := problem.String()
		if err := os.WriteFile(filepath, []byte(formatted), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing formatted puzzle: %v\n", err)
			os.Exit(1)
		}

		// Display the formatted puzzle
		fmt.Println(problem)
		fmt.Printf("\nFormatted puzzle written to %s\n", filepath)
	case "solve":
		fmt.Println("Problem:")
		fmt.Println(problem)

		// Solve the puzzle
		solution, err := problem.solve()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error solving problem: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nSolution:")
		fmt.Println(solution)

		// If output file is specified, write the solution to it
		if len(os.Args) == 4 {
			outputPath := os.Args[3]
			if err := os.WriteFile(outputPath, []byte(solution.String()), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing solution: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("\nSolution written to %s\n", outputPath)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintf(os.Stderr, "Valid commands are 'display', 'format', and 'solve'\n")
		os.Exit(1)
	}
}
