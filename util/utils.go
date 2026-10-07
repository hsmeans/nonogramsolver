package util

func GetColumn[T any](matrix [][]T, colIndex int) []T {
	column := make([]T, len(matrix))
	for i, row := range matrix {
		column[i] = row[colIndex]
	}
	return column
}
