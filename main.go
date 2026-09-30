package main

import (
	"fmt"
	"os"

	"github.com/hsmeans/nonogramsolver/util"
)

func main() {
	fmt.Println("Starting Nonogram Solver")

	if len(os.Args) < 2 {
		fmt.Println("Must provide a filename to a puzzle")
	}

	fpath := os.Args[1]

	_, err := util.LoadNonogramPuzzle(fpath)

	if err != nil {
		fmt.Println(err)
		return
	}

}
