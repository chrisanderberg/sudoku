package main

import (
	"strings"
	"testing"
)

func TestFromStringValidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  [81]int
	}{
		{
			name: "simple grid",
			input: `
				53..7....
				6..195...
				.98....6.
				8...6...3
				4..8.3..1
				7...2...6
				.6....28.
				...419..5
				....8..79`,
			want: [81]int{
				5, 3, 0, 0, 7, 0, 0, 0, 0,
				6, 0, 0, 1, 9, 5, 0, 0, 0,
				0, 9, 8, 0, 0, 0, 0, 6, 0,
				8, 0, 0, 0, 6, 0, 0, 0, 3,
				4, 0, 0, 8, 0, 3, 0, 0, 1,
				7, 0, 0, 0, 2, 0, 0, 0, 6,
				0, 6, 0, 0, 0, 0, 2, 8, 0,
				0, 0, 0, 4, 1, 9, 0, 0, 5,
				0, 0, 0, 0, 8, 0, 0, 7, 9,
			},
		},
		{
			name:  "single line input",
			input: "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79",
			want: [81]int{
				5, 3, 0, 0, 7, 0, 0, 0, 0,
				6, 0, 0, 1, 9, 5, 0, 0, 0,
				0, 9, 8, 0, 0, 0, 0, 6, 0,
				8, 0, 0, 0, 6, 0, 0, 0, 3,
				4, 0, 0, 8, 0, 3, 0, 0, 1,
				7, 0, 0, 0, 2, 0, 0, 0, 6,
				0, 6, 0, 0, 0, 0, 2, 8, 0,
				0, 0, 0, 4, 1, 9, 0, 0, 5,
				0, 0, 0, 0, 8, 0, 0, 7, 9,
			},
		},
		{
			name: "pretty printed input",
			input: `
				┌───────┬───────┬───────┐
				│ 5 3 · │ · 7 · │ · · · │
				│ 6 · · │ 1 9 5 │ · · · │
				│ · 9 8 │ · · · │ · 6 · │
				├───────┼───────┼───────┤
				│ 8 · · │ · 6 · │ · · 3 │
				│ 4 · · │ 8 · 3 │ · · 1 │
				│ 7 · · │ · 2 · │ · · 6 │
				├───────┼───────┼───────┤
				│ · 6 · │ · · · │ 2 8 · │
				│ · · · │ 4 1 9 │ · · 5 │
				│ · · · │ · 8 · │ · 7 9 │
				└───────┴───────┴───────┘`,
			want: [81]int{
				5, 3, 0, 0, 7, 0, 0, 0, 0,
				6, 0, 0, 1, 9, 5, 0, 0, 0,
				0, 9, 8, 0, 0, 0, 0, 6, 0,
				8, 0, 0, 0, 6, 0, 0, 0, 3,
				4, 0, 0, 8, 0, 3, 0, 0, 1,
				7, 0, 0, 0, 2, 0, 0, 0, 6,
				0, 6, 0, 0, 0, 0, 2, 8, 0,
				0, 0, 0, 4, 1, 9, 0, 0, 5,
				0, 0, 0, 0, 8, 0, 0, 7, 9,
			},
		},
		{
			name: "zeros instead of dots",
			input: `
				530070000
				600195000
				098000060
				800060003
				400803001
				700020006
				060000280
				000419005
				000080079`,
			want: [81]int{
				5, 3, 0, 0, 7, 0, 0, 0, 0,
				6, 0, 0, 1, 9, 5, 0, 0, 0,
				0, 9, 8, 0, 0, 0, 0, 6, 0,
				8, 0, 0, 0, 6, 0, 0, 0, 3,
				4, 0, 0, 8, 0, 3, 0, 0, 1,
				7, 0, 0, 0, 2, 0, 0, 0, 6,
				0, 6, 0, 0, 0, 0, 2, 8, 0,
				0, 0, 0, 4, 1, 9, 0, 0, 5,
				0, 0, 0, 0, 8, 0, 0, 7, 9,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromString(tt.input)
			if err != nil {
				t.Fatalf("fromString() error = %v", err)
			}
			if got.cells != tt.want {
				t.Errorf("fromString() = %v, want %v", got.cells, tt.want)
			}
		})
	}
}

func TestFromStringInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "too few cells",
			input:   "123456789",
			wantErr: "input contains only 9 cells, expected 81",
		},
		{
			name:    "too many cells",
			input:   strings.Repeat("1", 82),
			wantErr: "input contains more than 81 cells",
		},
		{
			name:    "invalid character",
			input:   "12345678x" + strings.Repeat("0", 72),
			wantErr: "invalid character 'x' at position 8",
		},
		{
			name: "duplicate in row",
			input: `
				123456781
				000000000
				000000000
				000000000
				000000000
				000000000
				000000000
				000000000
				000000000`,
			wantErr: "invalid sudoku problem: row 1 contains duplicate value 1",
		},
		{
			name: "duplicate in column",
			input: `
				100000000
				200000000
				300000000
				400000000
				500000000
				600000000
				700000000
				800000000
				100000000`,
			wantErr: "invalid sudoku problem: column 1 contains duplicate value 1",
		},
		{
			name: "duplicate in box",
			input: `
				123000000
				456000000
				781000000
				000000000
				000000000
				000000000
				000000000
				000000000
				000000000`,
			wantErr: "invalid sudoku problem: 3x3 box at position (1,1) contains duplicate value 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fromString(tt.input)
			if err == nil {
				t.Fatal("fromString() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("fromString() error = %v, want error containing %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringRoundTrip(t *testing.T) {
	input := `
		53..7....
		6..195...
		.98....6.
		8...6...3
		4..8.3..1
		7...2...6
		.6....28.
		...419..5
		....8..79`

	// Parse the input into a problem
	problem, err := fromString(input)
	if err != nil {
		t.Fatalf("fromString() error = %v", err)
	}

	// Convert problem to string and parse it back
	roundTrip, err := fromString(problem.String())
	if err != nil {
		t.Fatalf("fromString() error on round trip = %v", err)
	}

	// Compare the cells
	if problem.cells != roundTrip.cells {
		t.Errorf("Round trip failed: original = %v, got = %v", problem.cells, roundTrip.cells)
	}
}
