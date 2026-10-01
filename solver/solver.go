package solver

import (
	"slices"

	"github.com/hsmeans/nonogramsolver/util"
)

type Solution struct {
	Board [][]byte `json:"board"`
}

func initializeBoard(puzzle util.Puzzle) [][]byte {
	r, c := len(puzzle.Rows), len(puzzle.Columns)
	board := make([][]byte, r)
	for i := range board {
		board[i] = make([]byte, c)
	}

	return board
}

func SolveNonogram(puzzle util.Puzzle) {
	board := initializeBoard(puzzle)

	// while not solved (any spot is 0)
	for boardNotFilled(board) {
		// mark row hints
		markRows(puzzle, board)
		// mark col hints
		markCols(puzzle, board)
		// cross row hints
		crossRows(puzzle, board)
		// cross col hints
		crossCols(puzzle, board)
		// TODO: Need to fill
	}
}

func crossCols(puzzle util.Puzzle, board [][]byte) {
	for c, hints := range puzzle.Columns {
		// no hints, nothing to check
		if len(hints) == 0 {
			continue
		}

		for _, r := range crossLine(util.GetColumn(board, c), hints) {
			board[r][c] = 'x'
		}
	}
}

func crossRows(puzzle util.Puzzle, board [][]byte) {
	for r, hints := range puzzle.Rows {
		// no hints, nothing to check
		if len(hints) == 0 {
			continue
		}

		row := board[r]
		for _, c := range crossLine(row, hints) {
			row[c] = 'x'
		}
	}
}

func crossLine(line []byte, hints []int) []int {
	// count the marked sections of the line
	// if the number of marked sections equals the number of hints,
	// then we have known ranges that are the length of the sections +/- the size of the respective hint minus the length of the section
	// known range = len(section) +- (hint - len(section))
	// therefore, x out everything outside of the known range
	// if the sections equal the hints, the line will be filled out

	numSections := 0
	inSection := false
	for _, cell := range line {
		if cell == 'o' {
			if !inSection {
				inSection = true
				numSections++
			}
		} else {
			inSection = false
		}
	}

	if numSections == len(hints) {
		i, j := 0, 0

	}

	// If not all ranges are known ranges (num sections != num hints)
	// then check if there are known ranges within the uncomfirmed edges
	// i, hint[j] + 1
	// i - (hint[j] + 1), i
	// If a section exists and len(section) == respective hint:
	// cross out from start/end to the first cell after the section

	return []int{}
}

func markCols(puzzle util.Puzzle, board [][]byte) {
	for c, hints := range puzzle.Columns {
		// no hints, nothing to check
		if len(hints) == 0 {
			continue
		}

		for _, r := range markLine(util.GetColumn(board, c), hints) {
			board[r][c] = 'o'
		}
	}
}

func markRows(puzzle util.Puzzle, board [][]byte) {
	for r, hints := range puzzle.Rows {
		// no hints, nothing to check
		if len(hints) == 0 {
			continue
		}

		row := board[r]
		for _, c := range markLine(row, hints) {
			row[c] = 'o'
		}
	}
}

// Given a single row or column and its hints, return the indices of the cells
// that must be filled no matter how the hints end up being placed.
func markLine(line []byte, hints []int) []int {
	// go through the hints, capture the ranges for each check
	// for example, [7, 1] with no crosses will end up being [[[0, 6], [1, 7]], [[8, 8], [9, 9]]]
	hintRanges := make([][][]int, len(hints))
	for i := range hintRanges {
		hintRanges[i] = make([][]int, 0)
	}
	// check how it fits if starting at the rightmost possible spot
	// go through each spot in the row, and move through it and current hint, if you hit an X, reset the hint and continue
	// when a hint has been iterated, capture the range and move on to the next hint

	hintIdx := 0
	hint := hints[hintIdx]
	cur := hint

	low, high := 0, 0
	for i := 0; i < len(line); i++ {
		cell := line[i]
		// if there's an x, we cant put the hint here, restart on the next cell, otherwise keep going
		if cell != 'x' {
			high++
			cur--
		} else {
			low, high = i+1, i+1
			cur = hint
		}

		// if we've made it through the hint, save the range and move onto the next hint
		if cur == 0 {
			// high is one past the last cell of the run, so the inclusive range is [low, high-1]
			hintRange := [2]int{low, high - 1}
			hintRanges[hintIdx] = append(hintRanges[hintIdx], hintRange[:])
			hintIdx++
			// no more hints, move on
			if hintIdx >= len(hints) {
				break
			}
			hint = hints[hintIdx]
			cur = hint
			// runs need a gap between them, so skip the cell right after this run
			i++
			low, high = i+1, i+1
		}
	}

	// check how it fits if starting at the leftmost possible spot
	hintIdx = len(hints) - 1
	hint = hints[hintIdx]
	cur = hint

	end := len(line) - 1
	low, high = end, end
	for i := end; i >= 0; i-- {
		cell := line[i]

		// if there's an x, we cant put the hint here, restart on the next cell, otherwise keep going
		if cell != 'x' {
			low--
			cur--
		} else {
			low, high = i-1, i-1
			cur = hint
		}

		// if we've made it through the hint, save the range and move onto the next hint
		if cur == 0 {
			// low is one before the first cell of the run, so the inclusive range is [low+1, high]
			hintRange := [2]int{low + 1, high}
			hintRanges[hintIdx] = append(hintRanges[hintIdx], hintRange[:])
			hintIdx--
			// no more hints, move on
			if hintIdx < 0 {
				break
			}
			hint = hints[hintIdx]
			cur = hint
			// runs need a gap between them, so skip the cell right before this run
			i--
			low, high = i-1, i-1
		}
	}

	// the overlap of the two placements is what we can safely mark
	marked := make([]int, 0)
	for _, ranges := range hintRanges {
		lower, higher := ranges[0], ranges[1]
		marked = append(marked, getOverlap(lower, higher)...)
	}

	return marked
}

// Given two ranges ([1, 6], [4, 8])
// return a list of all numbers in the intersection of the ranges
// [4, 5, 6]
func getOverlap(lower, higher []int) []int {
	if lower[1] < higher[0] {
		return []int{}
	}

	x, y := max(lower[0], higher[0]), min(lower[1], higher[1])

	return util.RangeInts(x, y)
}

func boardNotFilled(board [][]byte) bool {
	for _, col := range board {
		if slices.Contains(col, 0) {
			return false
		}
	}

	return true
}
