package main

import (
	"fmt"
	"os"

	"github.com/hsmeans/nonogramsolver/solver"
	"github.com/hsmeans/nonogramsolver/util"
)

func main() {
	fmt.Println("Starting Nonogram Solver")

	if len(os.Args) < 2 {
		fmt.Println("Must provide a filename to a puzzle")
		return
	}

	fpath := os.Args[1]

	puzzle, err := util.LoadNonogramPuzzle(fpath)

	if err != nil {
		fmt.Println(err)
		return
	}

	solution, solved := solver.SolveNonogram(puzzle)
	if !solved {
		fmt.Println("No solution, partial board:")
	}
	fmt.Print(solution)
}
