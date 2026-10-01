package util

func RangeInts(x, y int) []int {
	if x > y {
		return []int{}
	}

	s := make([]int, 0, y-x+1)
	for i := x; i <= y; i++ {
		s = append(s, i)
	}
	return s
}

func GetColumn[T any](matrix [][]T, colIndex int) []T {
	column := make([]T, len(matrix))
	for i, row := range matrix {
		column[i] = row[colIndex]
	}
	return column
}
