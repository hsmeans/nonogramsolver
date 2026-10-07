package solver

import "strings"

type Solution struct {
	Board [][]byte `json:"board"`
}

func (s Solution) String() string {
	var sb strings.Builder
	for _, row := range s.Board {
		for _, cell := range row {
			switch cell {
			case 'o':
				sb.WriteByte('#')
			case 'x':
				sb.WriteByte('.')
			default:
				sb.WriteByte('?')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
