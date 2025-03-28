# sudoku
Sudoku solver written in Go

## Input Format
The program accepts Sudoku puzzles in two formats:

### Simple Format
A simple format using dots (.) or zeros (0) for empty cells, with optional whitespace:

### Pretty Format
A formatted grid using Unicode box-drawing characters and dots (·) for empty cells:

The program will output both the input puzzle and its solution in the pretty format.

## Usage
```bash
sudoku <command> <puzzle-file> [output-file]
```

### Commands
- `display`: Show the formatted puzzle
- `format`: Format the puzzle file in-place and display it
- `solve`: Solve the puzzle and display both problem and solution
  - Optionally write solution to output-file if specified

### Examples
```bash
# Display a formatted puzzle
$ sudoku display puzzle.txt

# Format a puzzle file in-place
$ sudoku format puzzle.txt

# Solve a puzzle and display the result
$ sudoku solve puzzle.txt

# Solve a puzzle and save the solution to a file
$ sudoku solve puzzle.txt solution.txt
```

## Example
```bash
$ cat puzzle.txt
53..7....
6..195...
.98....6.
8...6...3
4..8.3..1
7...2...6
.6....28.
...419..5
....8..79

$ sudoku solve puzzle.txt
Problem:
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
└───────┴───────┴───────┘

Solution:
┌───────┬───────┬───────┐
│ 5 3 4 │ 6 7 8 │ 9 1 2 │
│ 6 7 2 │ 1 9 5 │ 3 4 8 │
│ 1 9 8 │ 3 4 2 │ 5 6 7 │
├───────┼───────┼───────┤
│ 8 5 9 │ 7 6 1 │ 4 2 3 │
│ 4 2 6 │ 8 5 3 │ 7 9 1 │
│ 7 1 3 │ 9 2 4 │ 8 5 6 │
├───────┼───────┼───────┤
│ 9 6 1 │ 5 3 7 │ 2 8 4 │
│ 2 8 7 │ 4 1 9 │ 6 3 5 │
│ 3 4 5 │ 2 8 6 │ 1 7 9 │
└───────┴───────┴───────┘
```

## Algorithm
The solver works by converting the Sudoku puzzle into an exact cover problem:
- Each cell must contain exactly one number (81 constraints)
- Each row must contain each number exactly once (81 constraints)
- Each column must contain each number exactly once (81 constraints)
- Each 3x3 box must contain each number exactly once (81 constraints)

The resulting exact cover problem is solved using Knuth's Dancing Links algorithm, which is an efficient implementation of Algorithm X for solving exact cover problems.

## Error Handling
The program validates input puzzles and reports errors for:
- Invalid file paths
- Malformed puzzles (wrong size or invalid characters)
- Invalid puzzles (duplicate numbers in rows, columns, or boxes)
- Unsolvable puzzles

## License
MIT License. See [LICENSE](LICENSE) file for details.
