package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type Puzzle struct {
	Rows    [][]int `json:"rows"`
	Columns [][]int `json:"columns"`
}

func LoadNonogramPuzzle(fname string) (Puzzle, error) {
	fmt.Printf("Loading nonogram puzzle JSON at: %v\n", fname)
	content, err := os.ReadFile(fname)
	if err != nil {
		return Puzzle{}, fmt.Errorf("Error reading file %v: %v", fname, err)

	}

	var puzzle Puzzle
	err = json.Unmarshal(content, &puzzle)
	if err != nil {
		return Puzzle{}, fmt.Errorf("Error parsing JSON puzzle: %v", err)
	}

	return puzzle, nil
}
