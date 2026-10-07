package solver

import (
	"slices"

	"github.com/hsmeans/nonogramsolver/util"
)

func SolveNonogram(puzzle util.Puzzle) (Solution, bool) {
	board := initializeBoard(puzzle)
	solved := solve(puzzle, board)

	return Solution{Board: board}, solved
}

// Creates an empty 2d byte array based on the size of the puzzle
func initializeBoard(puzzle util.Puzzle) [][]byte {
	r, c := len(puzzle.Rows), len(puzzle.Columns)
	board := make([][]byte, r)
	for i := range board {
		board[i] = make([]byte, c)
	}

	return board
}

// Attempts to solve the puzzle. If it cannot be solved by traditional technique,
// then it will fall into trial and error. If it cannot be solved after that,
// a partial solution will be stored in the board parameter, and the function will return false.
func solve(puzzle util.Puzzle, board [][]byte) bool {
	// On the initial call, this will likely be the only part that will run.
	// If cannot be solved, return false.
	if !canSolve(puzzle, board) {
		return false
	}

	r, c, found := firstUnknown(board)
	// If there are no unknowns, then the puzzle has been solved
	if !found {
		return true
	}

	// Unknowns were found after exhausting traditional techniques.
	// Time for trial and error
	for _, guess := range []byte{'o', 'x'} {
		attempt := cloneBoard(board)
		attempt[r][c] = guess
		if solve(puzzle, attempt) {
			for i := range board {
				copy(board[i], attempt[i])
			}
			return true
		}
	}

	// Still can't solve after trial and error.
	// If initial call will return false with a partial solution stored in board parameter.
	return false
}

// Attepts to solve the puzzle.
// If the puzzle is a True Nonogram, this function is all that is needed.
// If there are no logic errors in the puzzle, return true.
func canSolve(puzzle util.Puzzle, board [][]byte) bool {
	for {
		// If any line results in a contradiction/logic error, it'll return false immediately
		rowsChanged, ok := solveRows(puzzle, board)
		if !ok {
			return false
		}
		colsChanged, ok := solveCols(puzzle, board)
		if !ok {
			return false
		}

		// If nothing changed, we've done as much as we can.
		if !rowsChanged && !colsChanged {
			return true
		}
	}
}

// Attempts to solve to columns of the board
func solveCols(puzzle util.Puzzle, board [][]byte) (bool, bool) {
	changed := false
	for c, hints := range puzzle.Columns {
		fills, crosses, ok := solveLine(util.GetColumn(board, c), hints)
		if !ok {
			return changed, false
		}
		for _, r := range fills {
			board[r][c] = 'o'
		}
		for _, r := range crosses {
			board[r][c] = 'x'
		}
		changed = changed || len(fills) > 0 || len(crosses) > 0
	}
	return changed, true
}

// Attempts to solve for the rows of the board
func solveRows(puzzle util.Puzzle, board [][]byte) (bool, bool) {
	changed := false
	for r, hints := range puzzle.Rows {
		row := board[r]
		fills, crosses, ok := solveLine(row, hints)
		if !ok {
			return changed, false
		}
		for _, c := range fills {
			row[c] = 'o'
		}
		for _, c := range crosses {
			row[c] = 'x'
		}
		changed = changed || len(fills) > 0 || len(crosses) > 0
	}
	return changed, true
}

// Attempts to solve the given line with the hints
func solveLine(line []byte, hints []int) (fills, crosses []int, ok bool) {
	left, ok := leftmostStarts(line, hints)
	// No combination of positions worked with the hints
	if !ok {
		return nil, nil, false
	}
	// No need to check here since if leftmost doesn't work then rightmost doesn't work
	right, _ := rightmostStarts(line, hints)

	covered := make([]bool, len(line))
	for hintIdx, hint := range hints {
		// These are guaranteed overlaps of left and right hints
		for i := right[hintIdx]; i < left[hintIdx]+hint; i++ {
			if line[i] == 0 {
				fills = append(fills, i)
			}
		}

		// These are in the bounds of possibility for the respective hint
		for i := left[hintIdx]; i < right[hintIdx]+hint; i++ {
			covered[i] = true
		}
	}

	// Anything not covered is a cross
	for i, cell := range line {
		if !covered[i] && cell == 0 {
			crosses = append(crosses, i)
		}
	}

	return fills, crosses, true
}

// Returns the leftmost starts for each hint
func leftmostStarts(line []byte, hints []int) ([]int, bool) {
	hintStarts := make([]int, len(hints))

	// DFS helper function for returning the first possible combination of leftmost starts
	var place func(hintIdx, pos int) bool
	place = func(hintIdx, pos int) bool {
		// If only one hint, return
		if hintIdx == len(hints) {
			return pos >= len(line) || !slices.Contains(line[pos:], 'o')
		}

		hint := hints[hintIdx]
		for cell := pos; cell+hint <= len(line); cell++ {
			// if the hint fits at the current cell, track it and check the next hint
			if fits(line, cell, hint) {
				hintStarts[hintIdx] = cell
				if place(hintIdx+1, cell+hint+1) {
					// We found valid leftmost starts
					return true
				}
			}

			/*
				If there's a filled cell here and its not a valid position,
				the current hint cannot start here, either because it doesn't fit
				or because subsequent hints cant be placed.
				This is a dead end for this combination.
			*/
			if line[cell] == 'o' {
				break
			}
		}
		// Try another combination.
		// If this returns from the initial call, then no possible combination
		// works, and there is a contradiction in the puzzle or board.
		return false
	}

	return hintStarts, place(0, 0)
}

// Returns the rightmost starts for each hint
func rightmostStarts(line []byte, hints []int) ([]int, bool) {
	// Reverse the line so we can re-use the logic from leftmostStarts
	rev := slices.Clone(line)
	slices.Reverse(rev)
	revHints := slices.Clone(hints)
	slices.Reverse(revHints)

	revStarts, ok := leftmostStarts(rev, revHints)
	if !ok {
		return nil, false
	}

	starts := make([]int, len(hints))
	// Since we reversed the line earlier,
	// we have to reverse the hints along with their board position
	for hintIdx, hint := range hints {
		starts[hintIdx] = len(line) - revStarts[len(hints)-1-hintIdx] - hint
	}

	return starts, true
}

// Returns if a hint can fit in the section starting at cell
func fits(line []byte, cell, hint int) bool {
	// If any part of the section contains a cross
	if slices.Contains(line[cell:cell+hint], 'x') {
		return false
	}

	// Check edge case if the hint reaches the end of the line (return true)
	// Otherwise, return true as long as there is not a filled cell after the section.
	return cell+hint == len(line) || line[cell+hint] != 'o'
}

// Returns true if there is an empty cell in the board
func firstUnknown(board [][]byte) (int, int, bool) {
	for r, row := range board {
		if c := slices.Index(row, 0); c >= 0 {
			return r, c, true
		}
	}

	return 0, 0, false
}

// Clones the board for guessing
func cloneBoard(board [][]byte) [][]byte {
	clone := make([][]byte, len(board))
	for i, row := range board {
		clone[i] = slices.Clone(row)
	}
	return clone
}
